package ;

typedef Problem =
{
	hole: Array<Point>,
	epsilon: Int,
	figure: {
		edges: Array<Point>,
		vertices: Array<Point>,
	},
	bonuses: Array<{
		bonus   :BonusKind,
		problem :Int,
		position:Array<Int>,
	}>
}

typedef Point = Array<Int>;
