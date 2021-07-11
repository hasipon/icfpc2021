package ;
import haxe.ds.Option;

typedef Problem = {
	isGlobalist:Bool,
	isWallhack:Bool,
	epsilon: Int,
	hole: Array<Point>,
	figure: {
		edges: Array<Point>,
	},
	distances:Array<Float>,
	bonuses  :Array<Bonus>,
	breakALeg: Option<Int>,
}
typedef ProblemSource =
{
	hole: Array<Point>,
	epsilon: Int,
	figure: {
		edges: Array<Point>,
		vertices: Array<Point>,
	},
	bonuses :Array<Bonus>,
}
typedef Bonus = {
	bonus   :BonusKind,
	problem :Int,
	position:Point,
}
typedef Point = Array<Int>;
