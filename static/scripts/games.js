"use strict";

var colors = [
	"#16a085", 
	"#8e44ad", 
	"#c0392b"
];
var numbers = [6,8,9];

function game0(main, first, second) {
	if (playerIndex == 1) {
		playerValue = first;
	} else {
		playerValue = second;
	}

	var lines = game0generator(8, 25, main, first, second);
	var list = $('<ul>').addClass('game0');

	$.each(lines, function(i) {
	    var li = $('<li/>')
	    	.attr('style', 'height: 35px')
	        .addClass('no-list')
	        .appendTo(list);
	    var subList = $('<ul>').appendTo(li);

	    numbers = lines[i]
	    $.each(numbers, function(j) {
	    	var num = numbers[j];
	    	$('<li onclick="numClicked(' + i + ',' + j + ')"/>')	    		
    			.attr('style', 'cursor: pointer; display: inline; color: ' + num.color)
	    		.addClass('no-list')
	    		.text(num.val)
	    		.appendTo(subList);
	    });
	});
	return list[0].outerHTML;
}

function numClicked (line, col) {
	hideNum(line, col)
	var $elem = $('ul.game0 li:nth-child(' + (line+1) + ') ul li:nth-child(' + (col+1) + ')');

	var val = parseInt($elem.html());
	score(val == playerValue);
	notify('hideNum(' + line + ',' + col + ')');
}

function hideNum (line, col) {
	var $elem = $('ul.game0 li:nth-child(' + (line+1) + ') ul li:nth-child(' + (col+1) + ')')
	var style = 'visibility: hidden;' + $elem.attr('style');
	$elem.attr('style', style);
}

function score (success) {
	$elem = $('#player-' + playerIndex + '-score');
	var val = parseInt($elem.html())
	if (success) {
		$elem.html(val + 1);
	} else {
		$elem.html(val - 1);
	}
}

function notify (str) {
	var ev = {
		'score_player_1': $player1score.html(),
		'score_player_2': $player2score.html(),
		'message': str
	};

	socket.emit('event', JSON.stringify(ev));
}

function game0generator(linesNum, lineLen, main, first, second) {
	var mainColor = getRandomFromArray(colors);

	var lines = [];
	for (var i = 0; i < linesNum; i++) {
		var line = [];
		for (var j = 0; j < lineLen; j++) {
			var val = main
			var rnd = Math.random()
			if (rnd < 2 / lineLen) {
				val = first;
			} else if (rnd < 4 / lineLen) {
				val = second;
			}

			var number = {
				val: val,
				color: mainColor
			};
			line.push(number);
		}
		lines.push(line)
	}
	
	return lines;
};


function getRandomFromArray(arr) {
	return arr[Math.floor(Math.random() * arr.length)];
}