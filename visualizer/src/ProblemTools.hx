package ;
import haxe.ds.Map;

class ProblemTools 
{
	public static function checkEpsilon(problem:Problem, ad:Float, pd:Float):Bool 
	{
		var e = problem.epsilon;
		return if (ad < pd) -(1000000 * ad) <= (e - 1000000) * pd else (1000000 * ad) <= (e + 1000000) * pd;
	}
	
	public static function checkGlobalEpsilon(problem:Problem, answer:Array<Array<Int>>):Bool
	{
		var value = 0.0;
		var e = (problem.epsilon * problem.figure.edges.length) / 1000000;
		for (ei => edge in problem.figure.edges)
		{
			var ax = answer[edge[0]][0] - answer[edge[1]][0];
			var ay = answer[edge[0]][1] - answer[edge[1]][1];
			var ad = ax * ax + ay * ay;
			var pd = problem.distances[ei];
			value += Math.abs(ad / pd - 1);
		}
		return value <= e;
	}
	
	
	public static function dislike(
		problem:Problem,
		answer:Array<Array<Int>>
	):Float
	{
		var dislike = 0.0;
		for (hole in problem.hole)
		{
			var min = Math.POSITIVE_INFINITY;
			var hx = hole[0];
			var hy = hole[1];
			for (a in answer)
			{
				var dx = hx - a[0];
				var dy = hy - a[1];
				var value = dx * dx + dy * dy;
				if (value < min) { min = value; }
			}
			dislike += min;
		}
		return dislike;
	}
	public static function failCount(
		problem:Problem,
		answer:Array<Array<Int>>
	):Int
	{
		var failCount = 0;
		var h0 = problem.hole[problem.hole.length - 1];
		var failedPoint = -1;
		var failedEdge0 = -1;
		var failedEdge1 = -1;
		
		for (h1 in problem.hole)
		{
			for (edge in problem.figure.edges)
			{
				if (
					intersect(
						h0, 
						h1, 
						answer[edge[0]],
						answer[edge[1]] 
					)
				) {
					if (problem.isWallhack)
					{
						if (failedPoint != -1)
						{
							if (
								failedPoint == edge[0] || 
								failedPoint == edge[1]
							)
							{
								continue;
							}
						}
						else
						{
							if (failedEdge0 == -1)
							{
								failedEdge0 = edge[0];
								failedEdge1 = edge[1];
								continue;
							}
							else
							{
								if (failedEdge0 == edge[0] || failedEdge0 == edge[1])
								{
									failedPoint = failedEdge0;
									continue;
								}
								if (failedEdge1 == edge[0] || failedEdge1 == edge[1])
								{
									failedPoint = failedEdge1;
									continue;
								}
							}
						}
					}
					failCount += 4;
				}
			}
			h0 = h1;
		}
		
		if (problem.isGlobalist)
		{
			var value = 0.0;
			var e = (problem.epsilon * problem.figure.edges.length) / 1000000;
			for (ei => edge in problem.figure.edges)
			{
				var ax = answer[edge[0]][0] - answer[edge[1]][0];
				var ay = answer[edge[0]][1] - answer[edge[1]][1];
				var ad = ax * ax + ay * ay;
				var pd = problem.distances[ei];
				value += Math.abs(ad / pd - 1);
			}
			if (value > e) {
				failCount += Math.ceil((value - e) / problem.epsilon / 1000000) + 2;
			}
		}
		else
		{
			for (ei => edge in problem.figure.edges)
			{
				var ax = answer[edge[0]][0] - answer[edge[1]][0];
				var ay = answer[edge[0]][1] - answer[edge[1]][1];
				var ad = ax * ax + ay * ay;
				var pd = problem.distances[ei];
				
				if (!checkEpsilon(problem, ad, pd)) {
					failCount += 1;
				}
			}
		}
		
		return failCount;
	}
	
	
	public static function intersect(
		a:Array<Int>,
		b:Array<Int>,
		c:Array<Int>,
		d:Array<Int>
	):Bool
	{
		var ax = a[0], ay = a[1];
		var bx = b[0], by = b[1];
		var cx = c[0], cy = c[1];
		var dx = d[0], dy = d[1];
		var s = (ax - bx) * (cy - ay) - (ay - by) * (cx - ax);
		var t = (ax - bx) * (dy - ay) - (ay - by) * (dx - ax);
		if (s * t >= 0) return false;
		
		var s  = (cx - dx) * (ay - cy) - (cy - dy) * (ax - cx);
		var t  = (cx - dx) * (by - cy) - (cy - dy) * (bx - cx);
		if (s * t >= 0) return false;
		return true;
	}
	
	public static function eval(dislike:Float, fail:Int):Float
	{
		return fail * 200 + dislike + (fail / 5) * dislike;
	}
	
	public static function checkPoint(
		problem:Problem, point:Array<Int>):Bool
	{
		var x = point[0];
		var y = point[1];
		var count = 0;
		var h0 = problem.hole[problem.hole.length - 1];
		for (h1 in problem.hole)
		{
			var x0 = h0[0] - x;
			var y0 = h0[1] - y;
			var x1 = h1[0] - x;
			var y1 = h1[1] - y;
			
			var cv = x0 * x1 + y0 * y1;
			var sv = x0 * y1 - x1 * y0;
			
			if (sv == 0 && cv <= 0)
			{
				return true;
			}
			
			if (y0 < y1)
			{
				var tmp = x0;
				x0 = x1;
				x1 = tmp;
				tmp = y0;
				y0 = y1;
				y1 = tmp;
			}
				
			if (y1 <= 0 && 0 < y0)
			{
				var a = x0 * (y1 - y0);
				var b = y0 * (x1 - x0);
				if(b < a){
                    ++count;
                }
			}
			h0 = h1;
		}
		return  count % 2 != 0;
	}
}
