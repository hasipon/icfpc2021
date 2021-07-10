//
// Created by shiota on 7/10/21.
//

#include<iostream>
#include<cstdlib>
#include<vector>

#define REP(i, b, n) for (Int i = b; i < Int(n); i++)
#define rep(i, n) REP(i, 0, n)
#define FOR(e, o) for (auto &&e : o)

using namespace std;

using Int = long long;
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
  vector<Point> holes;
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

bool validate(Problem &problem, vp &nowV){
  vp origV = problem.figure.V;
  FOR(e, problem.figure.E) {
    Int i = e.first;
    Int j = e.second;
    Int origD = distance(origV[i], origV[j]);
    Int nowD = distance(nowV[i], nowV[j]);
    Int diff = max(origD, nowD) - min(origD, nowD);
    if (diff * 1000000 > problem.epsilon * origD) {
      cerr << "Edge between(" << i << "," << j
           << ") has an invalid length: original: " << origD << "pose: " << nowD
           << endl;
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

int main() {
  Problem p;
  p.input();
  validate(p, p.figure.V);
  return 0;
}