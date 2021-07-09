package ;

typedef Problem =
{
	hole: Array<Point>,
	epsilon: Int,
	figure: {
		edges: Array<Point>,
		vertices: Array<Point>,
	}
}

typedef Point = Array<Int>;
