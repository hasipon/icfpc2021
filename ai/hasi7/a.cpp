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
	set<int> used = {109,131,143,146,148,149,147,145,135,119,136,126,139,132,118,111,116,115,130,133,144,134,124,120,114,102,88,99,79,94,73,52,37,33,28,13,12,21,11,10,19,31,19,10,8,5,14,18,6,1,0,2,7,23,22,42,50,71,62,54,77,87,113,98,122,128};
	const int TH = 6;
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
		vector<int> a {109,131,143,146,148,149,147,145,135,119,136,126,139,132,118,111,116,115,130,133,144,134,124,120,114,102,88,99,79,94,73,52,37,33,28,13,12,21,11,10,19,31,19,10};
		for (int i = 0; i < (int)a.size(); ++ i) {
			res[a[i]] = H[(i+60-44+2)%H.size()];
		}
	}
	{
		vector<int> a {8,5,14,18,6,1,0,2,7,23,22,42,50,71,62,54,77,87,113,98,122,128};
		for (int i = 0; i < (int)a.size(); ++ i) {
			res[a[i]] = H[(i+89-22+2)%H.size()];
		}
	}
	{
		vector<int> a {25,9,3,4,15,16};
		for (int i = 0; i < (int)a.size(); ++ i) {
			res[a[i]] = H[(i+52-11+2)%H.size()];
		}
	}
	for (auto [x,y] : res) cout << "[" << x << "," << y << "],";
	cout << endl;
}
