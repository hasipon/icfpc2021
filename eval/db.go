package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"
)

var defaultDB SQLiteDB

type BonusKey string

type SQLiteDB struct {
	*sqlx.DB
}

const schema = `
CREATE TABLE IF NOT EXISTS m_problem_setting
(
    problem_id       text,
    use_bonus        text default '',
    unlock_bonus_key text default '',
    PRIMARY KEY (problem_id)
);

CREATE TABLE IF NOT EXISTS solution
(
    id           text,
    json         text default '',
    problem_id   text default '',
    valid        int,
    dislike      real,
    dislike_s    text default '',
    use_bonus    text default '',
    unlock_bonus text default '',
    eval_message text default '',
    created_at   timestamp,
    updated_at   timestamp,
    PRIMARY KEY (id)
);
`
const indexes = `
CREATE INDEX IF NOT EXISTS ON SOLUTION_DISLIKE account(dislike);
`

type ProblemSetting struct {
	ProblemId      string   `db:"problem_id" json:"problem_id,omitempty"`
	UseBonus       string   `db:"use_bonus" json:"use_bonus,omitempty"`               // この問題で使うボーナスの設定
	UnlockBonusKey BonusKey `db:"unlock_bonus_key" json:"unlock_bonus_key,omitempty"` // この問題でアンロックする予定のBonusKey
}

type Solution struct {
	ID          string    `db:"id" json:"id,omitempty"`
	Json        string    `db:"json" json:"json,omitempty"`
	ProblemID   string    `db:"problem_id" json:"problem_id,omitempty"`
	Valid       int       `db:"valid" json:"valid,omitempty"`
	Dislike     float64   `db:"dislike" json:"dislike,omitempty"`
	DislikeS    string    `db:"dislike_s" json:"dislike_s,omitempty"`
	UseBonus    string    `db:"use_bonus" json:"use_bonus,omitempty"`
	UnlockBonus string    `db:"unlock_bonus" json:"unlock_bonus,omitempty"`
	EvalMessage string    `db:"eval_message" json:"eval_message,omitempty"`
	CreatedAt   time.Time `db:"created_at" json:"created_at,omitempty"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at,omitempty"`
}

// GenBonusKey bonus_key (ボーナスを使う問題ID_ボーナス名) を作る
func GenBonusKey(problemID string, bonusName string) BonusKey {
	problemIDInt, err := strconv.Atoi(problemID)
	if err != nil {
		panic(err)
	}
	return BonusKey(fmt.Sprintf("%04d_%s", problemIDInt, bonusName))
}

func (ps *ProblemSetting) IsBonusUseOk(bonusName string) bool {
	return ps.UseBonus == bonusName
}

func (db SQLiteDB) Ok() bool {
	return db.DB != nil
}

func (db SQLiteDB) Init() error {
	_, err := db.Exec(schema + indexes)
	return err
}

func (db SQLiteDB) Migrate() error {
	ctx := context.Background()
	tables := []string{
		"solution",
		"m_problem_setting",
	}

	tx, err := db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelDefault})
	if err != nil {
		return errors.Wrap(err, "Begin failed")
	}

	_, err = tx.Exec(schema)
	if err != nil {
		_ = tx.Rollback()
		return errors.Wrap(err, "Begin failed")
	}

	for _, table := range tables {
		tmp := table + "_tmp"
		_, err = tx.Exec(`ALTER TABLE ` + table + ` RENAME TO ` + tmp)
		if err != nil {
			_ = tx.Rollback()
			return errors.Wrap(err, "ALTER TABLE failed")
		}
	}

	_, err = tx.Exec(schema + indexes)
	if err != nil {
		_ = tx.Rollback()
		return errors.Wrap(err, "failed to create new tables")
	}

	for _, table := range tables {
		tmp := table + "_tmp"
		rows, err := tx.Query(`SELECT * FROM ` + tmp + ` LIMIT 1`)
		if err != nil {
			_ = tx.Rollback()
			return errors.Wrap(err, "SELECT failed")
		}

		columns, err := rows.Columns()
		if err != nil {
			_ = tx.Rollback()
			return errors.Wrap(err, "Columns() failed")
		}
		_ = rows.Close()

		_, err = tx.Exec(`INSERT INTO ` + table + `(` + strings.Join(columns, ",") + `) SELECT * FROM ` + tmp)
		if err != nil {
			_ = tx.Rollback()
			return errors.Wrap(err, "INSERT failed")
		}

		_, err = tx.Exec(`DROP TABLE ` + tmp)
		if err != nil {
			_ = tx.Rollback()
			return errors.Wrap(err, "DROP TABLE failed")
		}
	}

	return tx.Commit()
}

func (db SQLiteDB) RegisterSolution(name, problemID string, poseBytes []byte) (*Solution, error) {
	useBonus := ""
	var pose Pose
	if err := json.Unmarshal(poseBytes, &pose); err != nil {
		log.Fatal("pose:", err)
	}
	if 0 < len(pose.Bonuses) {
		if 1 < len(pose.Bonuses) {
			log.Fatal("TOO MANY BONUSES")
		}
		useBonus = pose.Bonuses[0].Bonus
	}

	now := time.Now()
	solution := &Solution{
		ID:        problemID + "-" + name,
		ProblemID: problemID,
		Json:      string(poseBytes),
		UseBonus:  useBonus,
		CreatedAt: now,
		UpdatedAt: now,
	}
	_, err := db.NamedExec(`
INSERT INTO solution (
    id,
    json,
    problem_id,
    valid,
    dislike,
    dislike_s,
    use_bonus,
    unlock_bonus,
    created_at,
    updated_at
) VALUES (
    :id,
    :json,
    :problem_id,
    :valid,
    :dislike,
    :dislike_s,
    :use_bonus,
    :unlock_bonus,
    :created_at,
    :updated_at
)`, solution)
	if err != nil {
		return nil, err
	}
	return solution, nil
}

func (db SQLiteDB) FindNoEvalSolution() (*Solution, error) {
	s := new(Solution)
	err := db.QueryRowx("SELECT * FROM solution WHERE valid = 0 ORDER BY RANDOM() LIMIT 1").StructScan(s)
	return s, err
}

func (db SQLiteDB) UpdateSolutionEvalResult(solution *Solution, dislikeStr string, valid bool, msg string, obtainBonuses []BonusKey) error {
	dislike, err := strconv.ParseFloat(dislikeStr, 64)
	if err != nil {
		return err
	}
	solution.Dislike = dislike
	solution.DislikeS = dislikeStr

	if valid {
		solution.Valid = 1
	} else {
		solution.Valid = 9
	}

	if msg != "" {
		solution.EvalMessage = msg
	}

	var unlockBonuses []string
	for _, b := range obtainBonuses {
		unlockBonuses = append(unlockBonuses, fmt.Sprint(b))
	}
	sort.Slice(unlockBonuses, func(i, j int) bool {
		return unlockBonuses[i] < unlockBonuses[j]
	})
	solution.UnlockBonus = strings.Join(unlockBonuses, ",")
	solution.UpdatedAt = time.Now()

	_, err = db.NamedExec(`
UPDATE solution
SET
	valid = :valid,
    dislike = :dislike,
    dislike_s = :dislike_s,
	unlock_bonus = :unlock_bonus,
	eval_message = :eval_message,
	updated_at = :updated_at
WHERE id = :id`,
		solution)
	return err
}

func (db SQLiteDB) FindBestSolution(problemID string) (*Solution, error) {
	solution := new(Solution)
	setting, err := db.GetProblemSetting(problemID)
	if err == sql.ErrNoRows {
		// No setting
		err = db.QueryRowx(
			"SELECT * FROM solution WHERE problem_id = ? AND valid = 1 AND use_bonus = '' ORDER BY dislike ASC LIMIT 1",
			problemID).StructScan(solution)
		return solution, err
	}
	if err == nil {
		// Use setting
		err = db.QueryRowx(
			"SELECT * FROM solution WHERE problem_id = ? AND valid = 1 AND (use_bonus = '' OR use_bonus = ?) AND unlock_bonus = ? ORDER BY dislike ASC LIMIT 1",
			problemID, setting.UseBonus, setting.UnlockBonusKey).StructScan(solution)
		return solution, err
	}
	return nil, err
}

func (db SQLiteDB) GetProblemSetting(problemID string) (*ProblemSetting, error) {
	s := &ProblemSetting{}
	err := db.QueryRowx("SELECT * FROM m_problem_setting WHERE problem_id = ?", problemID).StructScan(s)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (db SQLiteDB) GetAllProblemIDsInSubmission() ([]string, error) {
	var problemIDs []string
	err := db.Select(&problemIDs, "SELECT DISTINCT problem_id FROM solution")
	return problemIDs, err
}
