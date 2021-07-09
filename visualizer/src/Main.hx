package;

import haxe.Http;
import haxe.Json;
import js.Browser;
import js.html.CanvasElement;
import js.html.Document;
import js.html.Element;
import js.html.Event;
import js.html.InputElement;
import js.html.SelectElement;
import js.html.TextAreaElement;
import js.lib.Math;
import pixi.core.Application;
import pixi.core.graphics.Graphics;
import pixi.core.math.Point;
import pixi.interaction.InteractionEvent;
import tweenxcore.color.RgbColor;
using tweenxcore.Tools;

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
	static var gridGraphics:Graphics;
	
	static var problemIndex:Int;
	static var answer:Array<Array<Int>>;
	static var problem:Problem;
	static var scale:Float;
	
	static var selectedPoint:Int;
	static var startX   :Int;
	static var startY   :Int;
	static var startPoint:Point;
	static var problemCombo:SelectElement;
	static var answerText:TextAreaElement;
	static var autoDown  :Bool;
	
	static function main() 
	{
		canvas = cast(Browser.document.getElementById("pixi"), CanvasElement);
		
		canvas.width  = 1400;
		canvas.height = 980;
		
		problemCombo = cast Browser.document.getElementById("problem_combo");
		answerText   = cast Browser.document.getElementById("answer_text");
		var autoButton   = cast Browser.document.getElementById("auto_button");
		
		problemCombo.addEventListener("change", selectProblem);
		answerText  .addEventListener("input", onChangeAnswer);
		autoButton  .addEventListener("mousedown", onAutoDown);
		
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
	
	static function onChangeAnswer():Void 
	{
		try
		{
			var a:Array<Array<Int>> = Json.parse(answerText.value).vertices;
			if (a.length != answer.length) throw "invalid point length";
			for (point in a)
			{
				if (point.length != 2) throw "invalid point length";
			}
			for (i in 0...a.length)
			{
				answer[i][0] = Math.round(a[i][0]);
				answer[i][1] = Math.round(a[i][1]);
			}
			drawAnswer();
		}
		catch(e)
		{
			trace(e);
		}
	}
	
	static function fetchProblem(index:Int)
	{
		var h = new Http("./problems/" + index);
		
		h.onData = function(d) {
			problems.push(Json.parse(d));
			if (index == 1)
			{
				start();
			}
			var element = Browser.document.createElement('option');
			element.setAttribute("value", "" + (index));
			element.innerHTML = "" + (index);
			problemCombo.appendChild(element);
			
			fetchProblem(index + 1);
		}
		h.onError = function(e) {
		}
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
		Browser.window.requestAnimationFrame(onEnterFrame);
	}
	static function onEnterFrame(f:Float):Void
	{
		if (autoDown)
		{
			for (i in 0...1000)
			{
				var count      = [for (_ in answer) 0];
				var velocities = [for (_ in answer)[0.0, 0.0]];
				var e = problem.epsilon / 1000000;
				var matched = true;
				for (edge in problem.figure.edges)
				{
					var ax = answer[edge[0]][0] - answer[edge[1]][0];
					var ay = answer[edge[0]][1] - answer[edge[1]][1];
					var ad = ax * ax + ay * ay;
					var px = problem.figure.vertices[edge[0]][0] - problem.figure.vertices[edge[1]][0];
					var py = problem.figure.vertices[edge[0]][1] - problem.figure.vertices[edge[1]][1];
					var pd = px * px + py * py;
					
					if (Math.abs(ad / pd - 1) <= e) 
					{
					}
					else 
					{
						count[edge[0]] += 1; 
						count[edge[1]] += 1; 
						
						var v = (Math.sqrt(ad) - Math.sqrt(pd)) / 3;
						var d = Math.atan2(ay, ax);
						velocities[edge[0]][0] -= v * Math.cos(d);
						velocities[edge[0]][1] -= v * Math.sin(d);
						velocities[edge[1]][0] += v * Math.cos(d);
						velocities[edge[1]][1] += v * Math.sin(d);
						matched = false;
					}
				}
				if (matched) { break; }
				for (i in 0...answer.length)
				{
					var v = velocities[i];
					var c = count[i];
					if (c != 0)
					{
						answer[i][0] = Math.round(answer[i][0] + (v[0] / c) + Math.random() - 0.5);
						answer[i][1] = Math.round(answer[i][1] + (v[1] / c) + Math.random() - 0.5);
					}
				}
			}
			drawAnswer();
			outputAnswer();
		}
		Browser.window.requestAnimationFrame(onEnterFrame);
	}
	static function selectProblem(e:Event):Void
	{
		readProblem(problemCombo.selectedIndex);
	}
	static function onAutoDown():Void
	{
		autoDown = true;
	}
	static function onMouseUp():Void
	{
		if (selectedPoint >= 0)
		{
			outputAnswer();
		}
		
		autoDown = false;
		selectedPoint = -1;
		selectGraphics.clear();
	}
	
	static function outputAnswer():Void 
	{
		var dislike = 0.0;
		for (hole in problem.hole)
		{
			var min = Math.POSITIVE_INFINITY;
			var hx = hole[0];
			var hy = hole[1];
			for (a in answer)
			{
				var dx = a[0] - hx;
				var dy = a[1] - hy;
				var value = dx * dx + dy * dy;
				if (value < min) { min = value; }
			}
			dislike += min;
		}
		
		Browser.document.getElementById("dislike").textContent = "" + dislike;
		answerText.value = Json.stringify({vertices:answer});
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
			selectGraphics.drawCircle(x, y, 3);
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
		outputAnswer();
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
			answerGraphics.drawCircle(x, y, 3);
			
			first = false;
		}
		answerGraphics.endFill();
		var e = problem.epsilon / 1000000;
		for (edge in problem.figure.edges)
		{
			var ax = answer[edge[0]][0] - answer[edge[1]][0];
			var ay = answer[edge[0]][1] - answer[edge[1]][1];
			var ad = ax * ax + ay * ay;
			var px = problem.figure.vertices[edge[0]][0] - problem.figure.vertices[edge[1]][0];
			var py = problem.figure.vertices[edge[0]][1] - problem.figure.vertices[edge[1]][1];
			var pd = px * px + py * py;
			
			answerGraphics.lineStyle(
				2,
				if (Math.abs(ad / pd - 1) <= e) 
				{
					0x00CC00;
				}
				else if (ad > pd) 
				{
					var rate = (ad / pd).inverseLerp(1, 4).clamp();
					var color = new RgbColor(
						rate.lerp(0.6, 0.9),
						rate.lerp(0.4, 0.0),
						0
					);
					color.toRgbInt();
				}
				else 
				{
					trace(ad, pd);
					var rate = (pd / ad).inverseLerp(1, 4).clamp();
					var color = new RgbColor(
						0,
						rate.lerp(0.4, 0.0),
						rate.lerp(0.6, 0.9)
					);
					color.toRgbInt();
				}
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
