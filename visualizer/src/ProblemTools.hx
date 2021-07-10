package ;

class ProblemTools 
{
	public static function checkEpsilonValue(problem:Problem, ad:Float, pd:Float):Bool 
	{
		var e = problem.epsilon;
		return if (ad < pd) -(1000000 * ad) <= (e - 1000000) * pd else (1000000 * ad) <= (e + 1000000) * pd;
	}
	public static function failCount(
		problem:Problem,
		answer:Array<Int>
	):Int
	{
	//	var h0 = problem.hole[];
	//	for (h1 in problem.hole)
	//	{
	//		h0 = h1;
	//	}
		return 0;
	}
	
}
