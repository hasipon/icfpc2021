#include <bits/stdc++.h>
using namespace std;
typedef long long ll;
int N, E, V, EPS;
vector<pair<int,int>> H, G, P;
vector<vector<double>> dist;
vector<vector<double>> hDist {
	{},
	{52},
	{58.5918083, 27},
	{84.5918083, 53, 26},
	{94.86832981, 64.48332963, 37.48332963, 27},
	{64.48332963, 59.93329626, 37.48332963, 58.5918083, 52},
	{37.48332963, 37.48332963, 26, 52, 58.5918083, 27},
	{27, 58.5918083, 52, 78, 84.5918083, 53, 26},
};
vector<int> res;
double Ratio = 0.82;
void f(int p) {
	if (p == N) {
		for (auto x : res) cout << x << ",";
		cout << endl;
		return;
	}
	for (int i = 0; i < V; ++ i) {
		for (int j = 0; j < p; ++ j) {
			if (dist[i][res[j]] < hDist[p][j] * Ratio) goto next;
		}
		res[p] = i;
		f(p+1);
		next:;
	}
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
	dist = vector<vector<double>>(V, vector<double>(V, 1e100));
	for (int i = 0; i < V; ++ i) dist[i][i] = 0;
	for (auto [u, v] : G) {
		ll dx0 = P[u].first-P[v].first;
		ll dy0 = P[u].second-P[v].second;
		ll d0 = dx0*dx0 + dy0*dy0;
		dist[u][v] = dist[v][u] = sqrt(d0);
	}
	for (int k = 0; k < V; ++ k) for (int i = 0; i < V; ++ i) for (int j = 0; j < V; ++ j) dist[i][j] = min(dist[i][j], dist[i][k]+dist[k][j]);
	res.resize(N);
	f(0);
	vector<int> A = {54,26,10,1,3,31,35,57};
	auto pp = P;
	for (int i = 0; i < N; ++ i) pp[A[i]] = H[i];
	for (auto x : pp) cout << "[" << x.first << "," << x.second << "],";
	cout << endl;
}
