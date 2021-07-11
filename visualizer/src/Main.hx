package;

import haxe.Http;
import haxe.Json;
import haxe.Resource;
import js.Browser;
import js.html.CanvasElement;
import js.html.Event;
import js.html.InputElement;
import js.html.KeyboardEvent;
import js.html.SelectElement;
import js.html.TextAreaElement;
import js.lib.Math;
import pixi.core.Application;
import pixi.core.graphics.Graphics;
import pixi.core.math.Point;
import pixi.core.math.shapes.Rectangle;
import pixi.interaction.InteractionEvent;
import tweenxcore.color.RgbColor;
using tweenxcore.Tools;
using ProblemTools;
import Problem.ProblemSource;
import haxe.ds.Option;

typedef AvailableBonus = {bonus:BonusKind, from:Int, element:InputElement};
class Main 
{
	static var canvas:CanvasElement;
	static var pixi:Application;
	static var problems:Array<ProblemSource>;
	static var availableBonuses:Array<AvailableBonus>;
	
	static var left  :Int;
	static var right :Int;
	static var top   :Int;
	static var bottom:Int;
	
	static var problemGraphics:Graphics;
	static var answerGraphics:Graphics;
	static var selectGraphics:Graphics;
	static var hintGraphics:Graphics;
	
	static var problemIndex:Int;
	static var answer:Array<Array<Int>>;
	static var problem:Problem;
	static var scale:Float;
	
	static var selectRect:Rectangle;
	static var selectedPoints:Array<Int>;
	static var startPoint:Point;
	static var startAnswers:Array<Point>;
	static var problemCombo:SelectElement;
	static var answerText:TextAreaElement;
	static var autoDown  :Bool;
	static var fitDown   :Bool;
	static var randomDown:Bool;
	static var requestCount:Int;
	static var bestEval:Float;
	static var bestAnswer:Array<Array<Int>>;
	
	static function main() 
	{
		canvas = cast(Browser.document.getElementById("pixi"), CanvasElement);
		
		canvas.width  = 1180;
		canvas.height = 980;
		
		
		problemCombo = cast Browser.document.getElementById("problem_combo");
		answerText   = cast Browser.document.getElementById("answer_text");
		Browser.document.getElementById("fit_button"            ).addEventListener("mousedown", () -> { fitDown = true; });
		Browser.document.getElementById("auto_button"           ).addEventListener("mousedown", () -> { autoDown = true; });
		Browser.document.getElementById("fit_auto_button"       ).addEventListener("mousedown", () -> { fitDown = autoDown = true; });
		
		Browser.document.getElementById("random_button"         ).addEventListener("mousedown", () -> { randomDown = true; });
		Browser.document.getElementById("random_auto_button"    ).addEventListener("mousedown", () -> { randomDown = autoDown = true; });
		Browser.document.getElementById("random_fit_auto_button").addEventListener("mousedown", () -> { randomDown = fitDown = autoDown = true; });
		Browser.document.getElementById("reload_button").addEventListener("mousedown", () -> { readProblem(problemIndex); });
		Browser.document.getElementById("best_button").addEventListener("mousedown", () -> {
			answer = [for (a in bestAnswer)[for (i in a) i]];
			drawAnswer();
		});
		
		problemCombo.addEventListener("change", selectProblem);
		answerText  .addEventListener("input", onChangeAnswer);
		canvas      .addEventListener("keydown", onKeyDown);
		
		pixi = new Application({
			view  :canvas,
			transparent: true,
			width: canvas.width,
			height: canvas.height,
			autoResize: true,
		});
		pixi.stage.interactive = true;
		problems = [];
		selectRect = null;
		fetchProblem();
		requestCount = 0;
	}
	
	static function onKeyDown(e:KeyboardEvent):Void
	{
			switch (e.keyCode)
			{
				case KeyboardEvent.DOM_VK_A:
					if (e.ctrlKey)
					{
						untyped selectedPoints.length = 0;
						selectRect = null;
						for (i in 0...answer.length) {
							selectedPoints.push(i);
						}
						e.preventDefault();
						drawSelectedPoints();
					}

				case KeyboardEvent.DOM_VK_LEFT:
					rotate(15);
					e.preventDefault();
					
				case KeyboardEvent.DOM_VK_RIGHT:
					rotate(-15);
					e.preventDefault();
					
				case KeyboardEvent.DOM_VK_UP:
					rotate(90);
					e.preventDefault();
					
				case KeyboardEvent.DOM_VK_DOWN:
					rotate(-90);
					e.preventDefault();
					
				case KeyboardEvent.DOM_VK_Z:
					var cx = Math.round(canvas.width  / 2 / scale + left);
					for (i in selectedPoints)
					{
						var a = answer[i];
						a[0] = cx + cx - a[0];
					}
					drawAnswer();
					drawSelectedPoints();
					e.preventDefault();
					
				case KeyboardEvent.DOM_VK_X:
					var cy = Math.round(canvas.height / 2 / scale + top);
					for (i in selectedPoints)
					{
						var a = answer[i];
						a[1] = cy + cy - a[1];
					}
					drawAnswer();
					drawSelectedPoints();
					e.preventDefault();
					
				case _:
			}
	}
	static function rotate(degree:Float):Void
	{
		var cx = canvas.width  / 2 / scale + left;
		var cy = canvas.height / 2 / scale + top ;
		for (i in selectedPoints)
		{
			var a = answer[i];
			var dx = a[0] - cx;
			var dy = a[1] - cy;
			var d = Math.sqrt(dx * dx + dy * dy);
			var r = degree / 180 * Math.PI + Math.atan2(dy, dx);
			a[0] = Math.round(cx + d * Math.cos(r));
			a[1] = Math.round(cy + d * Math.sin(r));
		}
		drawAnswer();
		drawSelectedPoints();
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
	
	static function fetchProblem()
	{
		problems = Json.parse(Resource.getString("problems"));
		for (index => problem in problems)
		{
			var element = Browser.document.createElement('option');
			element.setAttribute("value", "" + (index + 1));
			element.innerHTML = "" + (index + 1);
			problemCombo.appendChild(element);
		}
		
		start();
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
		
		hintGraphics = new Graphics();
		pixi.stage.addChild(hintGraphics);
		
		selectedPoints = [];
		startAnswers = [];
		readProblem(0);
		
		pixi.stage.on("mousedown", onMouseDown);
		pixi.stage.on("mousemove", onMouseMove);
		Browser.document.addEventListener("mouseup", onMouseUp);
		Browser.window.requestAnimationFrame(onEnterFrame);
	}
	static function onEnterFrame(f:Float):Void
	{
		if (fitDown || autoDown || randomDown)
		{
			var shouldFix = (cast Browser.document.getElementById("fix_checkbox")).checked;
			var fixedMap = new Map();
			if (shouldFix)
			{
				for (hole in problem.hole)
				{
					for (i => a in answer)
					{
						if (hole[0] == a[0] && hole[1] == a[1])
						{
							fixedMap[i] = true;
						}
					}
				}
			}
			if (randomDown)
			{
				for (i in 0...1)
				{
					for (hole in problem.hole)
					{
						var i = Std.random(answer.length);
						if (shouldFix && fixedMap[i]) { continue; }
						
						var a = answer[i];
						var dx = a[0] - hole[0];
						var dy = a[1] - hole[1];
						if (dx != 0 || dy !=0)
						{
							var v = Math.sqrt(dx * dx + dy * dy);
							var d = Math.atan2(dy, dx);
							a[0] = Math.round(a[0] - v * Math.cos(d) * Math.random() * Math.random() + Math.random() - 0.5);
							a[1] = Math.round(a[1] - v * Math.sin(d) * Math.random() * Math.random() + Math.random() - 0.5);
						}
					}
				}
			}
			for (_ in 0...5)
			{
				if (fitDown)
				{
					for (i in 0...1)
					{
						for (hole in problem.hole)
						{
							var min = Math.POSITIVE_INFINITY;
							var target = 0;
							for (i in 0...answer.length)
							{
								var a = answer[i];
								var dx = a[0] - hole[0];
								var dy = a[1] - hole[1];
								var d = dx * dx + dy * dy;
								if (
									d < min &&
									(d == 0 || d + 20 < min || Math.random() < 0.5)
								)
								{
									min = d;
									target = i;
								}
							}
							if (min > 0)
							{
								if (shouldFix && fixedMap[target]) { continue; }
								var v = Math.sqrt(min);
								var a = answer[target];
								var dx = a[0] - hole[0];
								var dy = a[1] - hole[1];
								var d = Math.atan2(dy, dx);
								a[0] = Math.round(a[0] - v * Math.cos(d) + Math.random() - 0.5);
								a[1] = Math.round(a[1] - v * Math.sin(d) + Math.random() - 0.5);
							}
						}
					}
				}	
				if (autoDown)
				{
					for (i in 0...50000)
					{
						if (i % 10 == 0) { updateBest(); } 
						if (problem.isGlobalist && problem.checkGlobalEpsilon(answer)) { break; }
						var count      = [for (_ in answer) 0];
						var velocities = [for (_ in answer)[0.0, 0.0]];
						var e = problem.epsilon;
						var matched = true;
						for (ei => edge in problem.figure.edges)
						{
							var ax = answer[edge[0]][0] - answer[edge[1]][0];
							var ay = answer[edge[0]][1] - answer[edge[1]][1];
							var ad = ax * ax + ay * ay;
							var pd = problem.distances[ei];
				
							if (
								if (problem.isGlobalist) ad != pd else !problem.checkEpsilon(ad, pd)
							) 
							{
								count[edge[0]] += 1; 
								count[edge[1]] += 1; 
								
								var v = (Math.sqrt(ad) - Math.sqrt(pd)) / 5;
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
							if (shouldFix && fixedMap[i]) { continue; }
							var v = velocities[i];
							var c = count[i];
							if (c != 0)
							{
								if (c == 1 && Math.random() < 0.1) continue;
	 							answer[i][0] = Math.round(answer[i][0] + (v[0] / (c + 1)) + (Math.random() - 0.5));
	 							answer[i][1] = Math.round(answer[i][1] + (v[1] / (c + 1)) + (Math.random() - 0.5));
							}
						}
						
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
	static function updateBest():Void
	{
		var dislike = ProblemTools.dislike(problem, answer);
		var fail = ProblemTools.failCount(problem, answer);
		var eval = ProblemTools.eval(dislike, fail);
		if (eval < bestEval)
		{
			bestEval = eval;
			bestAnswer = [for (a in answer) [for (i in a) i]];
		}
	}
	static function onMouseUp():Void
	{
		if (selectedPoints.length >= 0)
		{
			outputAnswer();
		}
		untyped selectedPoints.length = 0;
		startPoint = null;
		hintGraphics.clear();
		selectGraphics.clear();
		selectGraphics.beginFill(0xCC0000);
		if (selectRect != null)
		{
			var i = 0;
			if (selectRect.width < 0)
			{
				var rx = selectRect.x + selectRect.width;
				selectRect.x = rx;
				selectRect.width = -selectRect.width;
			}
			if (selectRect.height < 0)
			{
				var ry = selectRect.y + selectRect.height;
				selectRect.y = ry;
				selectRect.height = -selectRect.height;
			}
			for (point in answer)
			{
				var x = (point[0] - left) * scale;
				var y = (point[1] - top ) * scale;
				if (selectRect.contains(x, y)) 
				{
					selectGraphics.drawCircle(x, y, 3);
					selectedPoints.push(i);
				}
				i += 1;
			}
			selectRect = null;
		}
		
		autoDown = false;
		fitDown = false;
		randomDown = false;
	}
	static public function drawSelectedPoints():Void
	{
		selectGraphics.clear();
		selectGraphics.beginFill(0xCC0000);
		for (selectedPoint in selectedPoints)
		{
			var point = answer[selectedPoint];
			var x = (point[0] - left) * scale;
			var y = (point[1] - top ) * scale;
			selectGraphics.drawCircle(x, y, 3);
		}
	}
	
	static function outputAnswer():Void 
	{
		answerText.value = getAnswer();
	}
	static function getAnswer():String
	{
		var bonuses:Array<Dynamic> = [];
		for (bonus in availableBonuses)
		{
			if (bonus.element.checked)
			{
				var b:Dynamic = { bonus:bonus.bonus , problem:bonus.from }
				switch (bonus.bonus)
				{
					case GLOBALIST:
					case BREAK_A_LEG:
					case WALLHACK:
				}
				bonuses.push(b);
			}
		}
		return Json.stringify({vertices:answer, bonuses: bonuses});
	}
	static function updateScore():Void
	{
		updateBest();
		var dislike = ProblemTools.dislike(problem, answer);
		var fail = ProblemTools.failCount(problem, answer);
		var eval = ProblemTools.eval(dislike, fail);
		Browser.document.getElementById("dislike").textContent = "" + dislike; 
		Browser.document.getElementById("fail").textContent = "" + fail; 
		Browser.document.getElementById("eval").textContent = "" + eval; 
		Browser.document.getElementById("best").textContent = "" + bestEval; 
		
		requestValidate();
	}
	static function requestValidate():Void
	{
		requestCount += 1;
		var r = requestCount;
		var h = new Http("../eval/" + (problemIndex + 1));
		h.onData = function(d) {
			if (requestCount == r)
			{
				Browser.document.getElementById("response").textContent = d;
			}
		}
		h.onError = function(e) {}
		h.setPostData(getAnswer());
		h.request(true);
	}	
	static function onMouseDown(e:InteractionEvent):Void
	{
		var nearest = 500.0;
		var i = 0;
		var selectedPoint = -1;
		for (point in answer)
		{
			var x = (point[0] - left) * scale;
			var y = (point[1] - top ) * scale;
			var dx = x - e.data.global.x;
			var dy = y - e.data.global.y;
			var d = dx * dx + dy * dy;
			if (nearest > d) 
			{
				selectedPoint = i;
				nearest = d;
			}
			i += 1;
		}
		if (selectedPoint == -1)
		{
			untyped selectedPoints.length = 0;
		}
		else if (selectedPoints.indexOf(selectedPoint) == -1)
		{
			untyped selectedPoints.length = 0;
			selectedPoints.push(selectedPoint);
		}
		else
		{
		}
		if (selectedPoints.length >= 1)
		{
			selectGraphics.clear();
			selectGraphics.beginFill(0xCC0000);
			startPoint = new Point(e.data.global.x, e.data.global.y);
			untyped startAnswers.length = 0;
			for (selectedPoint in selectedPoints)
			{
				var point = answer[selectedPoint];
				var x = (point[0] - left) * scale;
				var y = (point[1] - top ) * scale;
				selectGraphics.drawCircle(x, y, 3);
				startAnswers.push(new Point(point[0], point[1]));
			}
		}
		else
		{
			selectRect = new Rectangle();
			selectRect.x = e.data.global.x;
			selectRect.y = e.data.global.y;
			selectRect.width  = 0;
			selectRect.height = 0;
		}
		hintGraphics.clear();
		if (selectedPoints.length == 1)
		{
			var selectedPoint = selectedPoints[0];
			var sx = answer[selectedPoint][0];
			var sy = answer[selectedPoint][1];
			var points:Map<Int, Int> = [];
			for (ei => edge in problem.figure.edges)
			{
				if (edge[0] == selectedPoint) points.set(ei, edge[1]);
				if (edge[1] == selectedPoint) points.set(ei, edge[0]);
			}
			var l = if (sx - 300 < left  ) left   else sx - 300;
			var r = if (right < sx + 300 ) right  else sx + 300;
			var t = if (sy - 300 < top   ) top    else sy - 300;
			var b = if (bottom < sy + 300) bottom else sy + 300;
			for (x in l...r)
			{
				for (y in t...b)
				{
					var fail = false;
					for (ei => point in points)
					{
						var ax = answer[point][0] - x;
						var ay = answer[point][1] - y;
						var ad = ax * ax + ay * ay;
						var pd = problem.distances[ei];
				
						if (!problem.checkEpsilon(ad, pd))
						{
							fail = true;
						}
					}
					if (!fail)
					{
						var x = (x - left) * scale;
						var y = (y - top ) * scale;
						hintGraphics.beginFill(0x9999FF);
						hintGraphics.drawCircle(x, y, 4);
					}
				}
			}
		}
	}
	static function onMouseMove(e:InteractionEvent):Void
	{
		if (startPoint != null)
		{
			for (i in 0...selectedPoints.length)
			{
				var dx = e.data.global.x - startPoint.x;
				var dy = e.data.global.y - startPoint.y;
				answer[selectedPoints[i]][0] = Math.round(startAnswers[i].x + dx / scale);
				answer[selectedPoints[i]][1] = Math.round(startAnswers[i].y + dy / scale);
				drawAnswer();
			}
		}
		if (selectRect != null)
		{
			selectRect.width  = e.data.global.x - selectRect.x;
			selectRect.height = e.data.global.y - selectRect.y;
		
			selectGraphics.clear();
			selectGraphics.lineStyle(2, 0xCC0000);
			selectGraphics.drawRect(selectRect.x, selectRect.y, selectRect.width, selectRect.height);
		}
	}
	
	static function readProblem(index:Int):Void
	{
		bestEval = Math.POSITIVE_INFINITY;
		untyped selectedPoints.length = 0;
		problemIndex = index;
		var source = problems[index];
		answer = [];
		for (point in source.figure.vertices)
		{
			answer.push([point[0], point[1]]);
		}
		var bonusElement = Browser.document.getElementById("bonus");
		bonusElement.innerHTML = "";
		availableBonuses = [];
		for (i => p in problems)
		{
			for (bonus in p.bonuses)
			{
				if (bonus.problem == problemIndex + 1)
				{
					var element:InputElement = cast Browser.document.createElement('input');
					element.setAttribute("type", "checkbox");
					element.setAttribute("id", "bonus" + availableBonuses.length);
					element.addEventListener("input", () -> {
						updateBonuses();
						drawAnswer();
						outputAnswer();
					});
					var label = Browser.document.createElement('label');
					label.setAttribute("for", "bonus" + availableBonuses.length);
					label.textContent = bonus.bonus + " from " + (i + 1);
					bonusElement.appendChild(element);
					bonusElement.appendChild(label);
					
					availableBonuses.push({
						bonus: bonus.bonus,
						from: i + 1,
						element: element
					});
				}
			}
		}
		updateBonuses();
		
		left = right = problem.hole[0][0];
		top = bottom = problem.hole[0][1];
		for (point in source.hole)
		{
			if (left   > point[0]) left   = point[0];
			if (right  < point[0]) right  = point[0];
			if (top    > point[1]) top    = point[1];
			if (bottom < point[1]) bottom = point[1];
		}
		for (point in source.figure.vertices)
		{
			if (left   > point[0]) left   = point[0];
			if (right  < point[0]) right  = point[0];
			if (top    > point[1]) top    = point[1];
			if (bottom < point[1]) bottom = point[1];
		}
		left   -= 12;
		right  += 12;
		top    -= 12;
		bottom += 12;
		
		var w = (right - left);
		var h = (bottom - top);
		var sw = canvas.width / w;
		var sh = canvas.height / h;
		scale = if (sw > sh) sh else sw;
		
		var first = true;
		problemGraphics.clear();
		problemGraphics.beginFill(0xEFEFEF);
		problemGraphics.lineStyle(1, 0x788888);
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
		for (hole in problem.hole)
		{
			var x = (hole[0] - left) * scale;
			var y = (hole[1] - top ) * scale;
			problemGraphics.beginFill(0x899999);
			problemGraphics.drawCircle(x, y, 4);
		}
		problemGraphics.endFill();
		
		for (bonus in problem.bonuses)
		{
			var color = switch (bonus.bonus)
			{
				case BonusKind.GLOBALIST  :0xFFFF00;
				case BonusKind.BREAK_A_LEG:0x0000FF;
				case BonusKind.WALLHACK   :0x0FF9900;
			}
			problemGraphics.beginFill(color);
			var x = (bonus.position[0] - left) * scale;
			var y = (bonus.position[1] - top ) * scale;
			problemGraphics.drawCircle(x, y, 6);
		}
		
		
		drawAnswer();
		outputAnswer();
	}
	static function updateBonuses():Void
	{
		var source:ProblemSource = problems[problemIndex];
		untyped answer.length = source.figure.vertices.length;
		problem = {
			hole: source.hole,
			epsilon: source.epsilon,
			figure: {
				edges: [for (e in source.figure.edges) e],
			},
			bonuses:source.bonuses,
			distances:[],
			breakALeg: Option.None,
			isGlobalist: false,
			isWallhack : false,
		};
		for (bonus in availableBonuses)
		{
			if (bonus.element.checked)
			{
				switch (bonus.bonus)
				{
					case BonusKind.GLOBALIST  : problem.isGlobalist = true;
					case BonusKind.BREAK_A_LEG:
					case BonusKind.WALLHACK   : problem.isWallhack = true;
				}
			}
		}
		for (edge in source.figure.edges) 
		{
			var px = source.figure.vertices[edge[0]][0] - source.figure.vertices[edge[1]][0];
			var py = source.figure.vertices[edge[0]][1] - source.figure.vertices[edge[1]][1];
			problem.distances.push(px * px + py * py);
		}
	}
	
	static function drawAnswer():Void
	{
		answerGraphics.clear();
		var e = problem.epsilon;
		for (ei => edge in problem.figure.edges)
		{
			var ax = answer[edge[0]][0] - answer[edge[1]][0];
			var ay = answer[edge[0]][1] - answer[edge[1]][1];
			var ad = ax * ax + ay * ay;
			var pd = problem.distances[ei];
			
			answerGraphics.lineStyle(
				2,
				if (problem.isGlobalist)
				{
					if (ad == pd) { 0x00CC00; }
					else if (ad > pd) 
					{
						var rate = (ad / pd).inverseLerp(1, 4).clamp();
						var color = new RgbColor(
							rate.lerp(0.5, 0.9),
							rate.lerp(0.5, 0.0),
							0
						);
						color.toRgbInt();
					}
					else 
					{
						var rate = (pd / ad).inverseLerp(1, 4).clamp();
						var color = new RgbColor(
							0,
							rate.lerp(0.5, 0.0),
							rate.lerp(0.5, 0.9)
						);
						color.toRgbInt();
					}
				}
				else
				{
					if (problem.checkEpsilon(ad, pd)) { 0x00CC00; }
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
						var rate = (pd / ad).inverseLerp(1, 4).clamp();
						var color = new RgbColor(
							0,
							rate.lerp(0.4, 0.0),
							rate.lerp(0.6, 0.9)
						);
						color.toRgbInt();
					}
				}
			);
			var x = (answer[edge[0]][0] - left) * scale;
			var y = (answer[edge[0]][1] - top ) * scale;
			answerGraphics.moveTo(x, y);
			
			var x = (answer[edge[1]][0] - left) * scale;
			var y = (answer[edge[1]][1] - top ) * scale;
			answerGraphics.lineTo(x, y);
		}
		
		var first = true;
		answerGraphics.beginFill(0x00CC00);
		for (point in answer)
		{
			var x = (point[0] - left) * scale;
			var y = (point[1] - top ) * scale;
			answerGraphics.drawCircle(x, y, 3);
			
			first = false;
		}
		answerGraphics.endFill();
		updateScore();
	}
}
