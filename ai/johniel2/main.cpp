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
  // if (ccw(a1, a2, b1) == CCW::OTHER && ccw(a1, a2, b2) == CCW::OTHER) return false;
  // if (ccw(b1, b2, a1) == CCW::OTHER && ccw(b1, b2, a2) == CCW::OTHER) return false;
  if (ccw(a1, a2, b1) == CCW::OTHER) return false;
  if (ccw(a1, a2, b2) == CCW::OTHER) return false;
  if (ccw(b1, b2, a1) == CCW::OTHER) return false;
  if (ccw(b1, b2, a2) == CCW::OTHER) return false;
  return (ccw(a1, a2, b1) * ccw(a1, a2, b2) <= 0 &&
          ccw(b1, b2, a1) * ccw(b1, b2, a2) <= 0);
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
    for (lli j = 0; j <= mx_w; ++j) {
      pair<lli, lli> p = make_pair(i, j);
      bool f = winding_number(hole, p);
      for (int a = 0; a < hole.size(); ++a) {
        int b = (a + 1) % hole.size();
        if (hole[a] == p) f = true;
        f = f || (ccw(hole[a], hole[b], p) == CCW::OTHER);
      }
      if (f) v.push_back(p);
    }
  }
  each (i, hole) {
    assert(count(v.begin(), v.end(), i));
  }
  random_shuffle(v.begin(), v.end());
  return v;
}

const int N = 3000;
lli g[N][N];
vec<pair<lli, lli>> hole;
vec<pair<lli, lli>> vertex;
vec<pair<int, int>> edge;

lli epsilon;
int hole_size, edge_size, vectex_size;

lli dist(pair<lli, lli> a, pair<lli, lli> b)
{
  lli x = a.first - b.first;
  lli y = a.second - b.second;
  return x * x + y * y;
}

bool dist_epsilon(lli origin, lli moved)
{
  return abs(origin - moved) * 1000000 <= epsilon * origin;
}

vec<pair<lli, lli>> dist_filter(vec<pair<lli, lli>> candidates, map<int, pair<lli, lli>> fixed, int target)
{
  vec<pair<lli, lli>> v;
  each (c, candidates) {
    bool ok = true;
    each (f, fixed) {
      if (g[f.first][target] == -1) continue;
      ok = ok && dist_epsilon(g[f.first][target], dist(f.second, c));
      int x = 0;
      for (int i = 0; i < hole.size(); ++i) {
        int j = (i + 1) % hole.size();
        ok = ok && !intersect(hole[i], hole[j], f.second, c);
        // x += intersect(hole[i], hole[j], f.second, c);
      }
      // ok = ok && (x <= 1);
      unless (ok) break;
    }
    if (ok) v.push_back(c);
  }
  return v;
}

vec<pair<lli, lli>> near_hole(void)
{
  const int D = 2;
  vec<pair<lli, lli>> v = hole;
  for (int k = 0; k < hole.size(); ++k) {
    for (int i = -D; i <= +D; ++i) {
      for (int j = -D; j <= +D; ++j) {
        unless (i || j) continue;
        pair<lli, lli> p = make_pair(hole[k].first + i, hole[k].second + j);
        v.push_back(p);
      }
    }
  }
  return v;
}

static int try_cnt = 0;
lli solve(vec<pair<lli, lli>> candidates, map<int, pair<lli, lli>> fixed, vec<int> ord)
{
  if (ord.empty()) {
    static int solved = 0;
    static char buff[50];
    sprintf(buff, "./out.%d.json", solved++);
    clog << buff << endl;
    ofstream fout(buff);
    fout << "{" << endl;
    fout << "  \"vertices\":" << endl;
    int cnt = 0;
    fout << "    [" << endl;
    each (i, fixed) {
      fout << "      [" << i.second.first << ", " << i.second.second << "]";
      if (cnt + 1 != fixed.size()) fout << ", " << endl;
      else fout << endl;
      ++cnt;
    }
    fout << "    ]" << endl;
    fout << "}" << endl;
    if (solved == 10) assert(false);
    return 1;
  }

  if (++try_cnt == 5000) throw "";
  const int target = ord.back();
  ord.pop_back();
  if (ord.size() % 5 == 0) clog << "target: " << target << endl;
  assert(0 <= target && target < vertex.size());

  if (fixed.count(target)) {
    return solve(candidates, fixed, ord);
  } else {
    vec<pair<lli, lli>> v;
    if (fixed.size() < 1) {
      v = dist_filter(hole, fixed, target);
    } else {
      v = dist_filter(candidates, fixed, target);
    }
    lli succ = 0;
    each (i, v) {
      fixed[target] = i;
      succ += solve(candidates, fixed, ord);
    }
    return succ;
  }
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

  vec<int> ord;
  {
    set<int> vis;
    function<void(int)> rec = [&] (int curr) {
      vis.insert(curr);
      ord.push_back(curr);
      for (int next = 0; next < vertex.size(); ++next) {
        if (vis.count(next)) continue;
        if (g[curr][next] == -1) continue;
        rec(next);
      }
      return ;
    };
    for (int i = 0; i < vertex.size(); ++i) {
      if (!vis.count(i)) rec(i);
    }
    assert(vis.size() == vertex.size());
    assert(ord.size() == vis.size());
    reverse(ord.begin(), ord.end());
    clog << ord << endl;
  }

  auto cs = list_candidates(hole);

  // lli succ = 0;
  // for (int i = 0; i < vertex.size(); ++i) {
  //   for (int j = 0; j < hole.size(); ++j) {
  //     map<int, pair<lli, lli>> fixed;
  //     fixed[i] = hole[j];
  //     succ += solve(cs, fixed, ord);
  //   }
  // }
  // assert(succ);

  // auto cs = list_candidates(hole);
  // lli succ = 0;
  // for (int i1 = 0; i1 < vertex.size(); ++i1) {
  //   for (int i2 = 0; i2 < vertex.size(); ++i2) {
  //     if (g[i1][i2] == -1) continue;
  //     for (int j1 = 0; j1 < hole.size(); ++j1) {
  //       for (int j2 = 0; j2 < hole.size(); ++j2) {
  //         map<int, pair<lli, lli>> fixed;
  //         fixed[i1] = hole[j2];
  //         fixed[i2] = hole[j2];
  //         lli a = dist(vertex[i1], vertex[i2]);
  //         lli b = dist(hole[i1], hole[i2]);
  //         if (dist_epsilon(a, b)) {
  //           succ += solve(cs, fixed, ord);
  //         }
  //       }
  //     }
  //   }
  // }
  // assert(succ);

  vec<int> vidx;
  for (int i = 0; i < vertex.size(); ++i) {
    vidx.push_back(i);
  }
  random_shuffle(vidx.begin(), vidx.end());
  lli succ = 0;
  // for (int i1 = 0; i1 < vertex.size(); ++i1) {
  //   for (int i2 = i1 + 1; i2 < vertex.size(); ++i2) {
  // for (int i3 = i1 + 1; i3 < vertex.size(); ++i3) {
  //   for (int i4 = i3 + 1; i4 < vertex.size(); ++i4) {
  for (int _i1 = 0; _i1 < vertex.size(); ++_i1) {
    for (int _i2 = _i1 + 1; _i2 < vertex.size(); ++_i2) {
      int i1 = vidx[_i1];
      int i2 = vidx[_i2];
      if (g[i1][i2] == -1) continue;
      for (int _i3 = _i1 + 1; _i3 < vertex.size(); ++_i3) {
        for (int _i4 = _i3 + 1; _i4 < vertex.size(); ++_i4) {
          int i3 = vidx[_i3];
          int i4 = vidx[_i4];
          if (g[i3][i4] == -1) continue;
          for (int j1 = 0; j1 < hole.size(); ++j1) {
            for (int j2 = 0; j2 < hole.size(); ++j2) {
              for (int j3 = 0; j3 < hole.size(); ++j3) {
                for (int j4 = 0; j4 < hole.size(); ++j4) {
                  lli a = dist(hole[j1], hole[j2]);
                  lli b = dist(hole[j3], hole[j4]);
                  if (g[i1][i2] != a) continue;
                  if (g[i3][i4] != b) continue;
                  map<int, pair<lli, lli>> fixed;
                  fixed[i1] = hole[j1];
                  fixed[i2] = hole[j2];
                  fixed[i3] = hole[j3];
                  fixed[i4] = hole[j4];
                  cout << fixed << endl;
                  try {
                    try_cnt = 0;
                    succ += solve(cs, fixed, ord);
                  } catch (const char* e) {
                  }
                }
              }
            }
          }
        }
      }
    }
  }
  assert(succ);

  // {
  //   map<int, pair<lli, lli>> fixed;
  //   assert(solve(cs, fixed, ord));
  // }
  // return 0;
}
