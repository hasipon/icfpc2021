package ;

class ProblemTools 
{
	public static function checkEpsilonValue(problem:Problem, ad:Float, pd:Float):Bool 
	{
		var e = problem.epsilon;
		return if (ad < pd) -(1000000 * ad) <= (e - 1000000) * pd else (1000000 * ad) <= (e + 1000000) * pd;
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
		answer:Array<Array<Int>>,
		isGrobalist
	):Int
	{
		var failCount = 0;
		var h0 = problem.hole[problem.hole.length - 1];
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
					failCount += 4;
				}
			}
			h0 = h1;
		}
		
		trace(isGrobalist);
		if (isGrobalist)
		{
			var value = 0.0;
			var e = (problem.epsilon * problem.figure.edges.length) / 1000000;
			for (edge in problem.figure.edges)
			{
				var ax = answer[edge[0]][0] - answer[edge[1]][0];
				var ay = answer[edge[0]][1] - answer[edge[1]][1];
				var ad = ax * ax + ay * ay;
				var px = problem.figure.vertices[edge[0]][0] - problem.figure.vertices[edge[1]][0];
				var py = problem.figure.vertices[edge[0]][1] - problem.figure.vertices[edge[1]][1];
				var pd = px * px + py * py;
				value += Math.abs(ad / pd - 1);
			}
			if (value > e) {
				failCount += 10;
			}
		}
		else
		{
			for (edge in problem.figure.edges)
			{
				var ax = answer[edge[0]][0] - answer[edge[1]][0];
				var ay = answer[edge[0]][1] - answer[edge[1]][1];
				var ad = ax * ax + ay * ay;
				var px = problem.figure.vertices[edge[0]][0] - problem.figure.vertices[edge[1]][0];
				var py = problem.figure.vertices[edge[0]][1] - problem.figure.vertices[edge[1]][1];
				var pd = px * px + py * py;
				
				if (!checkEpsilonValue(problem, ad, pd)) {
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
}
