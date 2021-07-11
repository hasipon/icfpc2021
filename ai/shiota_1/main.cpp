//
// Created by shiota on 7/10/21.
//

#include<iostream>
#include<cstdlib>
#include<vector>
#include<sstream>
#include<algorithm>
#include "lib.hpp"
#include<cmath>



void output(vp &pose) {
  cout << "{\"vertices\":[";
  bool first = true;
  rep(i, pose.size()){
    if(!first){
      cout <<",";
    }
    first = false;
    cout <<"["<<pose[i].x <<","<<pose[i].y<<"]";
  }
  cout << "]}";
  cout << endl;
}
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
      target.y++;
    }
  }
  if(validate(p, pose)){
    return pose;
  }
  return vp();
}

const int magic = 1000;

void drushUp(Problem &p, vp &pose, vector<bool> used, const vii &floyd, const vii &nextTo){
  vp origV = p.figure.V;
  rep(i, 1000){
    bool allOk = true;
    FOR( u, used){
      if(!u)allOk = false;
    }
    if(allOk)break;
    FOR(e, p.figure.E){
      if(used[e.first] && used[e.second])continue;
      if(!used[e.first] && !used[e.second])continue;
      int from = e.first, to = e.second;
      if(used[to])swap(from, to);
      for(int i = -magic; i<magic; i++){
        for(int j = -magic; j<magic; j++){
          Point newTo;
          newTo.x = pose[to].x + i;
          newTo.y = pose[to].y + j;
          bool ok = false;
          FOR(next, nextTo[to]){
            Int origD = distance(origV[from], origV[to]);
            Int nowD = distance(pose[from], newTo);
            Int diff = max(origD, nowD) - min(origD, nowD);
            if (diff * 1000000 > p.epsilon * origD) {
              goto NEXT;
            }
            ok = true;
          }
          if(ok){
            pose[to] = newTo;
            used[to] = true;
          }
          NEXT:;
        }
      }
    }
  }

  output(pose);
}

int maxi;

void dfs(Problem &p, int holeId, vp &pose, vector<bool> &used, const vii &floyd, const vii &nextTo, int last, int skipShareEdge){
  if(!pruneEps(p, pose, used, floyd)){
    return;
  }
  if(maxi == holeId && holeId >= p.holes.size()){
    maxi = max(maxi, max);
    drushUp(p, pose, used, floyd, nextTo);
    return;
  }
  if(last != -1){
    FOR(next, nextTo[last]){
      if(used[next])continue;
      used[next] = true;
      pose[next] = p.holes[holeId];
      dfs(p, holeId+1, pose, used, floyd, nextTo, next, skipShareEdge);
      used[next] = false;
    }
  }
  if(skipShareEdge>0 ) {
    rep(i, pose.size()) {
      if (used[i])
        continue;
      used[i] = true;
      pose[i] = p.holes[holeId];
      dfs(p, holeId + 1, pose, used, floyd, nextTo, i, skipShareEdge - 1);
      used[i] = false;
    }
  }
}




int main() {
  srand(0);
  Problem p;
  p.input();

  int VN = p.figure.V.size();

  vii floyd;
  vii nextTo(p.figure.V.size());

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
    nextTo[i].push_back(j);
    nextTo[j].push_back(i);
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
  rep(i, VN+1){
    cout << "skipShareEdge" << ' ' << i << endl;
    maxi = 0;
    dfs(p, 0, pose, used, floyd, nextTo, -1, i);
  }

  return 0;
}