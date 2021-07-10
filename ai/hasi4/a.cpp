#include <bits/stdc++.h>
using namespace std;
typedef long long ll;
const int EPS = 15533;
bool check(ll dx0, ll dy0, ll dx1, ll dy1) {
	ll d0 = dx0*dx0 + dy0*dy0;
	ll d1 = (dx1*dx1 + dy1*dy1) * 1000000;
	return (1000000-EPS)*d0 <= d1 && d1 <= (1000000+EPS)*d0;
}
int main() {
	set<pair<int,int>> p4;
	set<pair<int,int>> p3;
	for (int x4 = 0; x4 <= 60; ++ x4) for (int y4 = x4; y4 <= 60; ++ y4) if (check(36,3,x4,y4)) {
		for (int x3 = -60; x3 <= 60; ++ x3) for (int y3 = -60; y3 <= 60; ++ y3) if (check(15,31,x3,y3) && check(21,28,x3-x4,y3-y4)) {
			p4.insert({x4,y4});
			p3.insert({x3,y3});
			cout << "[4] " << x4 << " " << y4 << endl;
			cout << " [3] " << x3 << " " << y3 << endl;
		}
	}
	cout << "-" << endl;
	for (auto [x4,y4] : p4) {
		for (int x0 = -60; x0 <= 60; ++ x0) for (int y0 = -60; y0 <= 60; ++ y0) if (check(14,34,x0,y0) && check(22,31,x0-x4,y0-y4)) {
			cout << "[4] " << x4 << " " << y4 << endl;
			cout << " [0] " << x0 << " " << y0 << endl;
		}
	}
	cout << "-" << endl;
	for (auto [x3,y3] : p3) {
		for (int x2 = -60; x2 <= 60; ++ x2) for (int y2 = -60; y2 <= 60; ++ y2) if (check(34,2,x2,y2) && check(19,29,x2-x3,y2-y3)) {
			cout << "[3] " << x3 << " " << y3 << endl;
			cout << " [2] " << x2 << " " << y2 << endl;
		}
	}
}
