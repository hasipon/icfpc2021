//
// Created by shiota on 7/10/21.
//

#include<iostream>
#include<cstdlib>
#include<vector>
#include<sstream>

#include "lib.hpp"
#include<cmath>


vp improveOne(const Problem &p, vp pose, int holeId, bool x){
  Point h = p.holes[holeId];
  Int mini = distance(h, pose[0]);
  Point &target = pose[0];
  REP(i, 1, pose.size()){
    Int tmp = distance(h, pose[i]);
    if(mini > tmp){
      mini = tmp;
      target = pose[i];
    }
  }
  if(x){
    if(h.x ==target.x){
      return vp();
    }
    if(h.x < target.x){
      target.x--;
    }else{
      target.x++;
    }
  } else {
    if(h.y ==target.y){
      return vp();
    }
    if(h.y < target.y){
      target.y--;
    }else{
      target.x++;
    }
  }
  if(validate(p, pose)){
    return pose;
  }
  return vp();
}

void dfs(Problem &p, int holeId, vp &pose, vector<bool> &used, const vii &floyd){
  if(!pruneEps(p, pose, used, floyd)){
    return;
  }
  if(holeId == p.holes.size()){
    cout << "{\"vertices\":[";
    bool first = true;
    rep(i, used.size()){
      if(!first){
        cout <<",";
      }
      first = false;
      if(used[i]){
        cout <<"["<<pose[i].x <<","<<pose[i].y<<"]";
      }else {
        cout <<"["<<p.figure.V[i].x <<","<<p.figure.V[i].y<<"]";
      }
    }
    cout << "]}";
    cout << endl;
    return;
  }
  rep(i, pose.size()){
    if(used[i])continue;
    used[i] = true;
    pose[i] = p.holes[holeId];
    dfs(p, holeId+1, pose, used, floyd);
    used[i] = false;
  }
}




int main() {
  srand(0);
  Problem p;
  p.input();

  int VN = p.figure.V.size();

  vii floyd;
  rep(i, p.figure.V.size()){
    floyd.push_back(vi(p.figure.V.size(), (1LL<<60)));
  }
  rep(i, VN){
    floyd[i][i] = 0;
  }
  FOR(e, p.figure.E){
    Int i = e.first;
    Int j = e.second;
    floyd[i][j] = floyd[j][i] = sqrt(distance(p.figure.V[i], p.figure.V[j]));
  }
  rep(k, VN){
    rep(i, VN){
      rep(j, VN){
        floyd[i][j] = min(floyd[i][j], floyd[i][k] + floyd[k][j]);
      }
    }
  }
  vp pose = p.figure.V;
  vector<bool> used(VN, false);
  dfs(p, 0, pose, used, floyd);

  return 0;
}