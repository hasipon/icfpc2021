//
// Created by shiota on 7/10/21.
//

#include<iostream>
#include<cstdlib>
#include<vector>
#include<sstream>

#include "lib.hpp"

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


int main() {
  srand(0);
  Problem p;
  p.input();
  vp pose = readPose();
  cerr << "now: " << dislike(p, pose) <<endl;
  cerr << validate(p, pose) <<endl;
  rep(i, p.holes.size()){
    while(true){
      bool x = (rand() % 2)==1;
      vp newPose = improveOne(p, pose, i, x);
      if(newPose.empty()){
        newPose = improveOne(p, pose, i, !x);
      }
      if(newPose.empty()){
        break;
      }
      pose = newPose;
      cerr << "now: " << dislike(p, pose) <<endl;
    }
  }

  return 0;
}