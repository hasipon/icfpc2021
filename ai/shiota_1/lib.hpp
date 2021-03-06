//
// Created by shiota on 7/10/21.
//

#ifndef SHIOTA_1_LIB_HPP
#define SHIOTA_1_LIB_HPP

#endif // SHIOTA_1_LIB_HPP
#define REP(i, b, n) for (Int i = b; i < Int(n); i++)
#define rep(i, n) REP(i, 0, n)
#define FOR(e, o) for (auto &&e : o)

using namespace std;

using Int = long long;
using vi = vector<Int>;
using vii = vector<vi>;
using pii = pair<Int, Int>;

class Point {
public:
  Int x;
  Int y;
};

using vp = vector<Point>;
class Figure {
public:
  vector<pii> E;
  vp V;
};


class Problem {
public:
  vp holes;
  Int epsilon;
  Figure figure;
  void input(){
    Int N, E, V;
    cin >> N >> E >> V >> epsilon;

    holes = vp(N);
    FOR(h, holes)cin >> h.x >> h.y;

    figure.E= vector<pii>(E);
    FOR(e, figure.E)cin >> e.first >> e.second;

    figure.V = vp(V);
    FOR(v, figure.V)cin >> v.x >> v.y;
  }
};

bool insidePolygon(const Point &p, const vp &holes){
  Int N = holes.size();
  Int cnt = 0;
  Int x = p.x;
  Int y = p.y;
  rep(i, holes.size()){
    Int x0 = holes[i].x -x;
    Int y0 = holes[i].y -y;
    Int x1 = holes[(i+1)%N].x -x;
    Int y1 = holes[(i+1)%N].y -y;

    Int cv = x0 * x1 + y0*y1;
    Int sv = x0 * y1 - x1 * y0;
    if(sv == 0 && cv <= 0){
      return true;
    }

    if(y0 >= y1){
      swap(x0, x1);
      swap(y0, y1);
    }

    if(y0 <= 0 && 0 < y1 && x0 *(y1-y0) > y0 * (x1-x0)){
      cnt++;
    }
  }
  return (cnt %2 == 1);
}


bool intersect(Point p1, Point p2, Point p3, Point p4){
  Int tc1 = (p1.x - p2.x) * (p3.y - p1.y) + (p1.y - p2.y) * (p1.x - p3.x);
  Int tc2 = (p1.x - p2.x) * (p4.y - p1.y) + (p1.y - p2.y) * (p1.x - p4.x);
  Int td1 = (p3.x - p4.x) * (p1.y - p3.y) + (p3.y - p4.y) * (p3.x - p1.x);
  Int td2 = (p3.x - p4.x) * (p2.y - p3.y) + (p3.y - p4.y) * (p3.x - p2.x);
  return tc1*tc2<0 and td1*td2<0;
}

Int distance(Point a, Point b) {
  Int x =  a.x - b.x;
  Int y = a.y - b.y;
  return x*x + y*y;
}

bool pruneEps(const Problem &problem, vp &pose, vector<bool> &used, const vii &floyd){
  vp origV = problem.figure.V;
  rep(i, used.size()){
    if(!used[i])continue;
    REP(j, i+1, used.size()){
      if(!used[j])continue;
      {
        double origD = floyd[i][j] * floyd[i][j];
        Int nowD = distance(pose[i], pose[j]);
        if(nowD > origD){
          // TODO tuning
          if ((nowD - origD) * 1000000 > (double)problem.epsilon * origD) {
            return false;
          }
        }
      }
    }
  }
  FOR(e, problem.figure.E) {
    Int i = e.first;
    Int j = e.second;
    if(!used[i] )continue;
    if(!used[j] )continue;
    Int origD = distance(origV[i], origV[j]);
    Int nowD = distance(pose[i], pose[j]);
    Int diff = max(origD, nowD) - min(origD, nowD);
    if (diff * 1000000 > problem.epsilon * origD) {
      return false;
    }
  }
  return true;
}

bool validate(const Problem &problem, vp &nowV){
  vp origV = problem.figure.V;
  FOR(e, problem.figure.E) {
    Int i = e.first;
    Int j = e.second;
    Int origD = distance(origV[i], origV[j]);
    Int nowD = distance(nowV[i], nowV[j]);
    Int diff = max(origD, nowD) - min(origD, nowD);
    if (diff * 1000000 > problem.epsilon * origD) {
      cerr << "Edge between(" << i << "," << j
           << ") has an invalid length: original: " << origD << " pose: " << nowD
           << " " << endl;
      return false;
    }
    auto H = problem.holes;
    rep(ii, H.size() - 1) {
      Int jj = (ii + 1) % H.size();
      if (intersect(H[ii], H[jj], nowV[i], nowV[j])) {
        cerr << "Edge between(" << i << "," << j << ") intersects: hole(" << ii
             << "," << jj << ")" << endl;
        return false;
      }
    }
  }
  return true;
}

Int dislike(Problem &p, vp &v){
  Int ret = 0;
  FOR(h, p.holes){
    Int mini = distance(h, v[0]);
    REP(i, 1, v.size()){
      mini = min(mini, distance(h, v[i]));
    }
    ret += mini;
  }
  return ret;
}

vp readPose(){
  string s;
  {
    string tmp;
    while(cin >> tmp){
      s += tmp;
    }
  }

  FOR(c, s){
    if(!isdigit(c))c=' ';
  }
  stringstream ss(s);
  vp ret;
  Point p;
  while(ss >> p.x >> p.y){
    ret.push_back(p);
  }
  return ret;
}

