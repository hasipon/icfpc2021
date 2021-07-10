// author: @___Johniel
// github: https://github.com/johniel/

#include <bits/stdc++.h>

#define each(i, c) for (auto& i : c)
#define unless(cond) if (!(cond))

using namespace std;

template<typename P, typename Q> ostream& operator << (ostream& os, pair<P, Q> p) { os << "(" << p.first << "," << p.second << ")"; return os; }
template<typename P, typename Q> istream& operator >> (istream& is, pair<P, Q>& p) { is >> p.first >> p.second; return is; }
template<typename T> ostream& operator << (ostream& os, vector<T> v) { os << "("; for (auto& i: v) os << i << ","; os << ")"; return os; }
template<typename T> istream& operator >> (istream& is, vector<T>& v) { for (auto& i: v) is >> i; return is; }
template<typename T> ostream& operator << (ostream& os, set<T> s) { os << "#{"; for (auto& i: s) os << i << ","; os << "}"; return os; }
template<typename K, typename V> ostream& operator << (ostream& os, map<K, V> m) { os << "{"; for (auto& i: m) os << i << ","; os << "}"; return os; }

template<typename T> inline T setmax(T& a, T b) { return a = std::max(a, b); }
template<typename T> inline T setmin(T& a, T b) { return a = std::min(a, b); }

using lli = long long int;
using ull = unsigned long long;
using point = complex<lli>;
using str = string;
template<typename T> using vec = vector<T>;

constexpr array<int, 8> di({0, 1, -1, 0, 1, -1, 1, -1});
constexpr array<int, 8> dj({1, 0, 0, -1, 1, -1, -1, 1});
constexpr lli mod = 1e9 + 7;

constexpr double eps = 1e-7;

template<typename T> inline pair<T, T> operator - (pair<T, T> a, pair<T, T> b)
{
  return make_pair(a.first - b.first, a.second - b.second);
}

template<typename T> inline pair<T, T> operator + (pair<T, T> a, pair<T, T> b)
{
  return make_pair(a.first - b.first, a.second - b.second);
}

template<typename T> inline pair<T, T> operator * (pair<T, T> a, T x)
{
  return make_pair(a.first * x, a.second * x);
}

template<typename T> inline pair<T, T> operator + (pair<T, T> a, T x)
{
  return make_pair(a.first + x, a.second + x);
}

template<typename T> inline pair<T, T>& operator += (pair<T, T>& a, pair<T, T> b)
{
  a.first += b.first;
  a.second += b.second;
  return a;
}

template<typename T> inline pair<T, T>& operator -= (pair<T, T>& a, pair<T, T> b)
{
  a.first -= b.first;
  a.second -= b.second;
  return a;
}

double angle(pair<lli, lli> a, pair<lli, lli> b)
{
  double p = a.first * b.second - a.second * b.first;
  double q = a.first * b.first + a.second * b.second;
  return atan2(p, q);
}

bool winding_number(vec<pair<lli, lli>> hole, pair<lli, lli> p)
{
  double z = 0;
  for(int i = 0; i < hole.size(); ++i) {
    int j = (i + 1) % hole.size();
    auto a = hole[i] - p;
    auto b = hole[j] - p;
    z += angle(a, b);
  }
  return abs(2 * M_PI - z) < eps;
}

vec<pair<lli, lli>> list_candidates(vec<pair<lli, lli>> hole)
{
  lli mx_h = -1;
  lli mx_w = -1;
  each (i, hole) {
    setmax(mx_h, i.first);
    setmax(mx_w, i.second);
  }
  vec<pair<lli, lli>> v;
  for (lli i = 0; i <= mx_h; ++i) {
    for (lli j = 0; j <= mx_h; ++j) {
      pair<lli, lli> p = make_pair(i, j);
      if (winding_number(hole, p)) v.push_back(p);
    }
  }
  return v;
}

lli dist(pair<lli, lli> p, pair<lli, lli> q = make_pair(0LL, 0LL))
{
  lli x = p.first - q.first;
  lli y = p.second - q.second;
  return x * x + y * y;
}

pair<lli, lli> normal(pair<lli, lli>  p)
{
  return make_pair(p.second, -p.first);
}

lli dot(pair<lli, lli> a, pair<lli, lli> b)
{
  return (a.first * b.first + a.second * b.second);
}

lli cross(pair<lli, lli> a, pair<lli, lli> b)
{
  return (a.first * b.second - a.second * b.first);
}

lli norm(pair<lli, lli> p)
{
  return norm(complex<lli>(p.first, p.second));
}

// dir. a -> b -> c
namespace CCW{
  enum{ RIGHT = 1, LEFT = -1, FRONT = -2, BACK = +2, OTHER = 0 };
};
int ccw(pair<lli, lli> a, pair<lli, lli> b, pair<lli, lli> c)
{
  b -= a;
  c -= a;
  if (cross(b, c) < 0) return CCW::RIGHT;
  if (cross(b, c) > 0) return CCW::LEFT;
  if (dot(b, c) < 0) return CCW::BACK;
  if (norm(b) < norm(c)) return CCW::FRONT;
  return CCW::OTHER;
}

bool intersect(pair<lli, lli> a1, pair<lli, lli> a2, pair<lli, lli> b1, pair<lli, lli> b2)
{
  return (ccw(a1, a2, b1) * ccw(a1, a2, b2) <= 0 &&
          ccw(b1, b2, a1) * ccw(b1, b2, a2) <= 0);
}

const int N = 3000;
lli g[N][N];
vec<pair<lli, lli>> hole;
vec<pair<lli, lli>> vertex;
vec<pair<int, int>> edge;

lli epsilon;
int hole_size, edge_size, vectex_size;

vec<pair<lli, lli>> dist_filter(vec<pair<lli, lli>> candidates, map<int, pair<lli, lli>> fixed, int target)
{
  vec<pair<lli, lli>> v;
  each (c, candidates) {
    bool ok = true;
    each (f, fixed) {
      if (g[f.first][target] == -1) continue;
      ok = ok && (g[f.first][target] == dist(f.second, c));
      for (int i = 0; i < hole.size(); ++i) {
        int j = (i + 1) % hole.size();
        ok = ok && !intersect(hole[i], hole[j], f.second, c);
      }
      unless (ok) break;
    }
    if (ok) v.push_back(c);
  }
  return v;
}

bool solve(vec<pair<lli, lli>> candidates, map<int, pair<lli, lli>> fixed, int target)
{
  if (target % 5 == 0) clog << "target: " << target << endl;
  if (target == -1) {
    // each (i, fixed) cout << i << endl;
    cout << "{" << endl;
    cout << "  \"vertices\":" << endl;
    int cnt = 0;
    cout << "    [" << endl;
    each (i, fixed) {
      cout << "      [" << i.second.first << ", " << i.second.second << "]";
      if (cnt + 1 != fixed.size()) cout << ", " << endl;
      else cout << endl;
      ++cnt;
    }
    cout << "    ]" << endl;
    cout << "}" << endl;
    return true;
  }

  vec<pair<lli, lli>> v = dist_filter(candidates, fixed, target);
  // cout << target << ' ' << v << endl;
  each (i, v) {
    fixed[target] = i;
    if (solve(candidates, fixed, target - 1)) return true;
  }
  return false;
}

int main(int argc, char *argv[])
{
  ios_base::sync_with_stdio(0);
  cin.tie(0);
  cout.setf(ios_base::fixed);
  cout.precision(15);

  cin >> hole_size >> edge_size >> vectex_size >> epsilon;

  hole.resize(hole_size);
  vertex.resize(vectex_size);
  edge.resize(edge_size);
  cin >> hole >> edge >> vertex;

  fill(&g[0][0], &g[N - 1][N - 1] + 1, -1);
  each (e, edge) {
    g[e.first][e.second] = g[e.second][e.first] = dist(vertex[e.first], vertex[e.second]);
  }

  // TODO: epsilonのことも考えてあげて。
  // TODO: 次数順に並べる。
  // TODO: cs[i]とcs[i+1]の間には辺があった方が良い。
  auto cs = list_candidates(hole);
  // cout << v << ' ' << v.size() << endl;
  const bool succ = solve(cs, map<int, pair<lli, lli>>(), vertex.size() - 1);
  assert(succ);
  return 0;
}
