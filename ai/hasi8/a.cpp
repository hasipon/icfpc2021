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
	set<int> used = {115,107,101,100,111,104,93,95,88,86,75,73,79,83,91,92,84,72,77,69,63,60,62,65,59,53,52,50,48,31,17,19,18,26,14,9,4,0,3,2,1,5,8,15,21,20,13,22,32,30,41,44,34,38};
	const int TH = 8;
	for (int i = 0; i < 2*N; ++ i) {
		map<int,vector<vector<int>>> b;
		for (auto [u,v] : hoge[i%N]) {
			b[u].push_back({v,u});
			b[v].push_back({u,v});
			for (auto aa : a[u]) {
				if (find(aa.begin(), aa.end(), v) != aa.end()) continue;
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
				if (find(aa.begin(), aa.end(), u) != aa.end()) continue;
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
	auto res = vector<pair<int,int>>(V);
	{
		vector<int> a {115,107,101,100,111,104,93,95,88,86,75,73,79,83,91,92,84,72,77,69,63,60,62,65,59,53,52,50,48};
		for (int i = 0; i < (int)a.size(); ++ i) {
			res[a[i]] = H[(i+65-29+2)%H.size()];
		}
	}
	{
		vector<int> a {31,17,19,18,26,14,9,4,0,3,2,1,5,8,15};
		for (int i = 0; i < (int)a.size(); ++ i) {
			res[a[i]] = H[(i+85-15+2)%H.size()];
		}
	}
	{
		vector<int> a {21,20,13,22,32,30,41,44,34,38};
		for (int i = 0; i < (int)a.size(); ++ i) {
			res[a[i]] = H[(i+14-10+2)%H.size()];
		}
	}
	for (auto [x,y] : res) cout << "[" << x << "," << y << "],";
	cout << endl;
}
