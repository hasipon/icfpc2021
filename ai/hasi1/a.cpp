#include <bits/stdc++.h>
using namespace std;
typedef long long ll;
int N, E, V, EPS;
vector<pair<int,int>> H, G, P;
set<pair<int,int>> Es;
vector<int> vv;
bool check(int p, int j, int v0, int v1) {
	ll dx0 = P[v0].first - P[v1].first;
	ll dy0 = P[v0].second - P[v1].second;
	ll dx1 = H[p].first - H[j].first;
	ll dy1 = H[p].second - H[j].second;
	ll d0 = dx0*dx0 + dy0*dy0;
	ll d1 = (dx1*dx1 + dy1*dy1) * 1000000;
	return (1000000-EPS)*d0 <= d1 && d1 <= (1000000+EPS)*d0;
}
void f(int p, int q) {
	if (p == N) {
		vector<pair<int,int>> ans(V, {-1,-1});
		for (int i = 0; i < N; ++ i) {
			ans[vv[i]] = H[i];
		}
		for (int i = 0; i < V; ++ i) {
			if (i) cout << ",";
			cout << "[" << ans[i].first << "," << ans[i].second << "]";
		}
		cout << endl;
		return;
	}
	for (int i = 0; i < V; ++ i) if (((q>>i)&1)==0) {
		vv[p] = i;
		for (int j = 0; j < p; ++ j) if (Es.count({vv[p], vv[j]}) || Es.count({vv[j], vv[p]})) {
			if (!check(p, j, vv[p], vv[j])) goto next;
		}
		f(p+1, q|(1<<i));
		next:;
	}
}
int main() {
	cin >> N >> E >> V >> EPS;
	H.resize(N);
	G.resize(E);
	P.resize(V);
	for (auto& x : H) cin >> x.first >> x.second;
	for (auto& x : G) cin >> x.first >> x.second;
	for (auto& x : P) cin >> x.first >> x.second;
	Es = set<pair<int,int>>(G.begin(), G.end());
	vv.resize(N);
	f(0, 0);
}
