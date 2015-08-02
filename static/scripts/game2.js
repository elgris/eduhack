"use strict";

var game2remains = 0;
var leftSide = true;

function game2(first, second, invert, game2data) {
	leftSide = !invert;
	if (playerIndex == 1) {
		playerValue = first;
	} else {
		leftSide = !leftSide;
		playerValue = second;
	}

	var mainColor = getRandomFromArray(colors);

	var lines = game2data;
	var list = $('<ul>').addClass('game game2');

	$.each(lines, function(i) {
	    var li = $('<li/>')
	    	.attr('style', 'height: 35px')
	        .addClass('no-list')
	        .appendTo(list);
	    var subList = $('<ul>').appendTo(li);

	    var numbers = lines[i]
	    var hl = numbers.length / 2
	    $.each(numbers, function(j) {
	    	var num = numbers[j];
	    	if (num == playerValue) {
	    		game2remains++;
	    	}
	    	var correct = (j < hl) == leftSide ? 'true' : 'false'
	    	$('<li onclick="game2Clicked(' + i + ',' + j + ',' + correct + ')"/>')	    		
    			.attr('style', 'cursor: pointer; display: inline; color: ' + mainColor)
	    		.addClass('no-list')
	    		.text(num)
	    		.appendTo(subList);
	    	if (j < hl && (j + 1 >= hl)) {
	    		$('<li />')
    			.attr('style', 'cursor: pointer; display: inline; color: ' + mainColor)
	    		.addClass('no-list')	    		
	    		.text('|')
	    		.appendTo(subList);
	    	}
	    });
	});
	return list[0].outerHTML;
}

function game2Clicked (line, col, correct) {
	if (game2remains <= 0) {
		return;
	}
	hideTile(line, col)
	var $elem = $('ul.game2 li:nth-child(' + (line+1) + ') ul li:nth-child(' + (col+1) + ')');

	var val = parseInt($elem.html());
	if (val == playerValue) {
		game2remains--;
	}
	score((val == playerValue) && correct);
	notify('hideTile(' + line + ',' + col + ')');

	if (game2remains <= 0) {
		finish('You have finished your part!', true);
	}
}
