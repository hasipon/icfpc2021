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
	set<int> used = {91,87,88,81,75,73,60,74,84,78,64,46,31,21,9,11,8,15,16,6,0,1,10,23,26,34,55,53,79,80,82,68,57,51};
	const int TH = 5;
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
	for (int i = 0; i < V; ++ i) res[i] = {400,i*4};
	{
		vector<int> a {87,91,88,81,75,73,60,74,84,78,64,46,31};
		for (int i = 0; i < (int)a.size(); ++ i) {
			res[a[i]] = H[(i+40-a.size()+2)%H.size()];
		}
	}
	{
		vector<int> a {21,9,11,8,15,16,6,0,1,10,23};
		for (int i = 0; i < (int)a.size(); ++ i) {
			res[a[i]] = H[(i+51-a.size()+2)%H.size()];
		}
	}
	{
		vector<int> a {26,34,55,53,79,80,82,68,57,51};
		for (int i = 0; i < (int)a.size(); ++ i) {
			res[a[i]] = H[(i+22-a.size()+2)%H.size()];
		}
	}
	for (auto [x,y] : res) cout << "[" << x << "," << y << "],";
	cout << endl;
}
