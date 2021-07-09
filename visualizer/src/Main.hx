package;

import haxe.Http;
import haxe.Json;
import js.Browser;
import js.html.CanvasElement;
import pixi.core.Application;
import pixi.core.graphics.Graphics;
import pixi.core.math.Point;
import pixi.interaction.InteractionEvent;

class Main 
{
	static var canvas:CanvasElement;
	static var pixi:Application;
	static var problems:Array<Problem>;
	static var left  :Int;
	static var right :Int;
	static var top   :Int;
	static var bottom:Int;
	
	static var problemGraphics:Graphics;
	static var answerGraphics:Graphics;
	static var selectGraphics:Graphics;
	
	static var problemIndex:Int;
	static var answer:Array<Array<Int>>;
	static var problem:Problem;
	static var scale:Float;
	
	static var selectedPoint:Int;
	static var startX   :Int;
	static var startY   :Int;
	static var startPoint:Point;
	
	static function main() 
	{
		canvas = cast(Browser.document.getElementById("pixi"), CanvasElement);
		
		canvas.width  = 800;
		canvas.height = 700;
		
		pixi = new Application({
			view  :canvas,
			transparent: true,
			width: canvas.width,
			height: canvas.height,
			autoResize: true,
		});
		pixi.stage.interactive = true;
		problems = [];
		fetchProblem(1);
	}
	
	
	static function fetchProblem(index:Int)
	{
		var h = new Http("./problems/" + index);
		
		h.onData = function(d) {
			problems.push(Json.parse(d));
			fetchProblem(index + 1);
			if (index == 1)
			{
				start();
			}
		}
		h.onError = function(e) {}
		h.request();
	}
	
	static function start():Void
	{
		var background = new Graphics();
		background.beginFill(0xCCCCCC);
		background.drawRect(0, 0, canvas.width, canvas.height);
		pixi.stage.addChild(background     );
		
		problemGraphics = new Graphics();
		pixi.stage.addChild(problemGraphics);
		
		answerGraphics = new Graphics();
		pixi.stage.addChild(answerGraphics);
		
		selectGraphics = new Graphics();
		pixi.stage.addChild(selectGraphics);
		
		readProblem(0);
		
		startPoint = new Point();
		pixi.stage.on("mousedown", onMouseDown);
		pixi.stage.on("mousemove", onMouseMove);
		Browser.document.addEventListener("mouseup", onMouseUp);
	}
	
	static function onMouseUp():Void
	{
		selectedPoint = -1;
		selectGraphics.clear();
	}
	
	static function onMouseDown(e:InteractionEvent):Void
	{
		selectedPoint = -1;
		
		var nearest = 500;
		var i = 0;
		for (point in answer)
		{
			var x = (point[0] - left) * scale;
			var y = (point[1] - top ) * scale;
			var dx = x - e.data.global.x;
			var dy = y - e.data.global.y;
			var d = dx * dx + dy * dy;
			if (nearest > d) selectedPoint = i;
			i += 1;
		}
		if (selectedPoint >= 0)
		{
			selectGraphics.clear();
			selectGraphics.beginFill(0xCC0000);
			
			var point = answer[selectedPoint];
			var x = (point[0] - left) * scale;
			var y = (point[1] - top ) * scale;
			selectGraphics.drawCircle(x, y, 2);
			startPoint.x = e.data.global.x;
			startPoint.y = e.data.global.y;
			startX = point[0];
			startY = point[1];
		}
	}
	static function onMouseMove(e:InteractionEvent):Void
	{
		if (selectedPoint >= 0)
		{
			var dx = e.data.global.x - startPoint.x;
			var dy = e.data.global.y - startPoint.y;
			answer[selectedPoint][0] = Math.round(startX + dx / scale);
			answer[selectedPoint][1] = Math.round(startY + dy / scale);
			trace(selectedPoint, dx, dy, scale);
			drawAnswer();
		}
	}
	
	static function readProblem(index:Int):Void
	{
		selectedPoint = -1;
		problem = problems[index];
		left = right = problem.hole[0][0];
		top = bottom = problem.hole[0][1];
		problemIndex = index;
		
		for (point in problem.hole)
		{
			if (left   > point[0]) left   = point[0];
			if (right  < point[0]) right  = point[0];
			if (top    > point[1]) top    = point[1];
			if (bottom < point[1]) bottom = point[1];
		}
		for (point in problem.figure.vertices)
		{
			if (left   > point[0]) left   = point[0];
			if (right  < point[0]) right  = point[0];
			if (top    > point[1]) top    = point[1];
			if (bottom < point[1]) bottom = point[1];
		}
		left   -= 3;
		right  += 3;
		top    -= 3;
		bottom += 3;
		
		var w = (right - left);
		var h = (bottom - top);
		var sw = canvas.width / w;
		var sh = canvas.height / h;
		scale = if (sw > sh) sh else sw;
		
		var first = true;
		problemGraphics.clear();
		problemGraphics.beginFill(0xEFEFEF);
		problemGraphics.lineStyle(1, 0x88888);
		for (hole in problem.hole)
		{
			var x = (hole[0] - left) * scale;
			var y = (hole[1] - top ) * scale;
			if (first)
			{
				problemGraphics.moveTo(x, y);
			}
			else
			{
				problemGraphics.lineTo(x, y);
			}
			first = false;
		}
		problemGraphics.endFill();
		
		answer = [];
		for (point in problem.figure.vertices)
		{
			answer.push([point[0], point[1]]);
		}
		drawAnswer();
	}
	
	static function drawAnswer():Void
	{
		var first = true;
		answerGraphics.clear();
		answerGraphics.beginFill(0x00CC00);
		for (point in answer)
		{
			var x = (point[0] - left) * scale;
			var y = (point[1] - top ) * scale;
			answerGraphics.drawCircle(x, y, 2);
			
			first = false;
		}
		answerGraphics.endFill();
		var e = problem.epsilon / 1000000;
		for (edge in problem.figure.edges)
		{
			var ax = answer[edge[0]][0] - answer[edge[1]][0];
			var ay = answer[edge[0]][1] - answer[edge[1]][1];
			var ad = Math.sqrt(ax * ax + ay * ay);
			var px = problem.figure.vertices[edge[0]][0] - problem.figure.vertices[edge[1]][0];
			var py = problem.figure.vertices[edge[0]][1] - problem.figure.vertices[edge[1]][1];
			var pd = Math.sqrt(px * px + py * py);
			
			answerGraphics.lineStyle(
				1,
				if (Math.abs(ad - pd) < e) 0x00CC00 else if (ad > pd) 0xCC0000 else 0x0000CC
			);
			var x = (answer[edge[0]][0] - left) * scale;
			var y = (answer[edge[0]][1] - top ) * scale;
			answerGraphics.moveTo(x, y);
			
			var x = (answer[edge[1]][0] - left) * scale;
			var y = (answer[edge[1]][1] - top ) * scale;
			answerGraphics.lineTo(x, y);
		}
	}
}
