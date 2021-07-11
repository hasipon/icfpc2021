package ;

typedef Problem =
{
	hole: Array<Point>,
	epsilon: Int,
	figure: {
		edges: Array<Point>,
		vertices: Array<Point>,
	},
	bonuses: Array<Bonus>
}

typedef Bonus = {
	bonus   :BonusKind,
	problem :Int,
	position:Point,
}
typedef Point = Array<Int>;
