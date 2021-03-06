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
	set<int> used = {154,140,121,112,100,81,79,56,31,14,13,1,0,12,7,15,206,211,207,210,209,213,214,212,208,188,182,164};
	const int TH = 6;
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
	//for (int i = 0; i < V; ++ i) res[i] = {400,i*4};
	{
		vector<int> a {154,140,121,112,100,81,79,56,31,14,13,1,0,12,7,15};
		for (int i = 0; i < (int)a.size(); ++ i) {
			res[a[i]] = H[(i+67-a.size()+2)%H.size()];
		}
	}
	{
		vector<int> a {206,211,207,210,209,213,214,212,208,188,182,164};
		for (int i = 0; i < (int)a.size(); ++ i) {
			res[a[i]] = H[(i+49-a.size()+2)%H.size()];
		}
	}
	{
		vector<int> a {4,6,2,9,27,18};
		for (int i = 0; i < (int)a.size(); ++ i) {
			res[a[i]] = H[(i+18-a.size()+2)%H.size()];
		}
	}
	for (auto [x,y] : res) cout << "[" << x << "," << y << "],";
	cout << endl;
	/*
	vector<pair<int,int>> pose = {{0,247},{3,289},{12,194},{15,320},{40,264},{54,192},{56,304},{71,234},{75,189},{80,329},{92,352},{103,311},{105,288},{106,184},{108,244},{120,344},{123,322},{126,221},{128,293},{129,203},{142,240},{231,278},{158,274},{172,73},{179,170},{180,226},{181,331},{185,104},{188,3},{190,299},{191,390},{193,260},{194,111},{193,74},{203,359},{205,43},{213,316},{210,408},{213,346},{223,341},{225,203},{228,310},{228,0},{195,63},{230,190},{232,249},{234,363},{232,370},{238,144},{249,279},{240,394},{240,67},{244,110},{248,303},{245,362},{246,166},{248,334},{249,220},{283,226},{254,283},{254,204},{252,376},{259,364},{261,410},{262,440},{293,305},{250,254},{267,280},{267,131},{296,138},{278,365},{310,160},{269,18},{268,332},{241,276},{272,173},{275,199},{275,48},{275,345},{274,375},{262,375},{282,209},{249,234},{284,153},{215,318},{287,113},{218,290},{292,373},{290,180},{296,403},{299,282},{330,304},{302,179},{301,358},{301,445},{303,77},{284,244},{361,253},{303,327},{310,140},{313,181},{318,152},{331,223},{253,302},{332,260},{321,34},{367,291},{322,370},{322,421},{376,316},{324,357},{326,247},{374,295},{327,390},{312,206},{361,296},{331,157},{343,233},{345,232},{334,200},{384,297},{321,211},{338,269},{321,243},{321,245},{384,306},{338,237},{339,381},{383,302},{339,329},{377,314},{354,318},{343,130},{342,199},{342,328},{314,183},{351,281},{352,321},{345,382},{346,96},{347,214},{373,305},{349,126},{344,265},{355,314},{338,291},{368,223},{361,254},{354,408},{409,298},{333,318},{344,301},{360,304},{358,57},{358,171},{314,283},{354,299},{362,197},{363,324},{371,233},{363,220},{371,340},{372,371},{373,147},{379,214},{360,314},{327,313},{385,312},{376,273},{349,331},{330,340},{345,277},{396,216},{383,398},{389,269},{370,290},{386,252},{391,325},{393,229},{350,291},{393,352},{348,253},{367,284},{398,293},{345,337},{383,242},{404,218},{408,254},{358,305},{381,264},{414,276},{417,393},{418,323},{422,349},{426,231},{427,196},{428,368},{432,332},{438,207},{439,258},{446,223},{461,251}};
	ll dislike = 0;
	for (auto h : H) {
		ll s = 1<<30;
		for (auto p : pose) {
			s = min(s, d(h, p));
		}
		if (s > 0) cout << h.first << " " << h.second << " " << s << endl;
		dislike += s;
	}
	cout << dislike << endl;
	*/
}
