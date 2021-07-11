#include <bits/stdc++.h>
using namespace std;
typedef long long ll;
int N, E, V, EPS;
vector<pair<int,int>> H, G, P;
ll d(pair<int,int> p0, pair<int,int> p1) {
	ll dx = p0.first - p1.first;
	ll dy = p0.second - p1.second;
	return dx*dx + dy*dy;
}
int main() {
	cin >> N >> E >> V >> EPS;
	H.resize(N);
	G.resize(E);
	P.resize(V);
	for (auto& x : H) {
		cin >> x.first >> x.second;
	}
	for (auto& x : G) cin >> x.first >> x.second;
	for (auto& x : P) cin >> x.first >> x.second;
	vector<ll> Hd(N);
	for (int i = 0; i < N; ++ i) {
		Hd[i] = d(H[i], H[(i+1)%N]);
	}
	vector<vector<pair<int,int>>> hoge(N);
	for (int i = 0; i < E; ++ i) {
		ll d1 = d(P[G[i].first], P[G[i].second]);
		for (int j = 0; j < N; ++ j) {
			ll d0 = Hd[j];
			if ((1000000-EPS)*d0 <= d1*1000000 && d1*1000000 <= (1000000+EPS)*d0) {
				hoge[j].push_back(G[i]);
			}
		}
	}
	map<int,int> counter;
	map<int,vector<vector<int>>> a;
	set<int> used = {52,55,70,85,80,98,102,107,114,120,116,125,130,128,130,131,129,123,122,111,93,104,29,23,13,8,3,9,2,0,1,4,6,5,7,12,11,16,21,34,42,28,46,};
	const int TH = 11;
	for (int i = 0; i < 2*N; ++ i) {
		map<int,vector<vector<int>>> b;
		for (auto [u,v] : hoge[i%N]) {
			b[u].push_back({v,u});
			b[v].push_back({u,v});
			for (auto aa : a[u]) {
				aa.push_back(v);
				b[v].push_back(aa);
				bool ok = true;
				for (int x : aa) if (used.count(x)) { ok = false; break; }
				if (ok) {
					++ counter[(int)aa.size()];
					if (aa.size() >= TH) {
						cout << i << " " << aa.size() << endl;
						for (auto x : aa) cout << x << ",";
						cout << endl;
					}
				}
			}
			for (auto aa : a[v]) {
				aa.push_back(u);
				b[u].push_back(aa);
				bool ok = true;
				for (int x : aa) if (used.count(x)) { ok = false; break; }
				if (ok) {
					++ counter[(int)aa.size()];
					if (aa.size() >= TH) {
						cout << i << " " << aa.size() << endl;
						for (auto x : aa) cout << x << ",";
						cout << endl;
					}
				}
			}
		}
		a = b;
	}
	for (auto [x, y] : counter) cout << x << " " << y << endl;
	auto res = P;
	{
		vector<int> a {52,55,70,85,80,98,102,107,114,120,116,125,124,128,130,131,129,123,122,111,93,79};
		for (int i = 0; i < (int)a.size(); ++ i) {
			res[a[i]] = H[(i+34-22+2)%H.size()];
		}
	}
	{
		vector<int> a {29,23,13,8,3,9,2,0,1,4,6,5,7,12,11,16,21,34,42,28,46,56};
		for (int i = 0; i < (int)a.size(); ++ i) {
			res[a[i]] = H[(i+85-22+2)%H.size()];
		}
	}
	{
		vector<int> a {86,103,101,83,86,96,88,76,74,69,53};
		for (int i = 0; i < (int)a.size(); ++ i) {
			res[a[i]] = H[(i+52-11+2)%H.size()];
		}
	}
	for (auto [x,y] : res) cout << "[" << x << "," << y << "],";
	cout << endl;
}
