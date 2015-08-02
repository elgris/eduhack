"use strict";

var colors = [
	"#16a085", 
	"#8e44ad", 
	"#c0392b"
];

var game0remains = 0;

function game0(first, second, game0data) {
	if (playerIndex == 1) {
		playerValue = first;
	} else {
		playerValue = second;
	}

	var mainColor = getRandomFromArray(colors);

	var lines = game0data;
	var list = $('<ul>').addClass('game game0');

	$.each(lines, function(i) {
	    var li = $('<li/>')
	    	.attr('style', 'height: 35px')
	        .addClass('no-list')
	        .appendTo(list);
	    var subList = $('<ul>').appendTo(li);

	    var numbers = lines[i]
	    $.each(numbers, function(j) {
	    	var num = numbers[j];
	    	if (num == playerValue) {
	    		game0remains++;
	    	}
	    	$('<li onclick="game0Clicked(' + i + ',' + j + ')"/>')	    		
    			.attr('style', 'cursor: pointer; display: inline; color: ' + mainColor)
	    		.addClass('no-list')
	    		.text(num)
	    		.appendTo(subList);
	    });
	});
	return list[0].outerHTML;
}

function game0Clicked (line, col) {
	if (game0remains <= 0) {
		return;
	}
	hideTile(line, col)
	var $elem = $('ul.game0 li:nth-child(' + (line+1) + ') ul li:nth-child(' + (col+1) + ')');

	var val = parseInt($elem.html());
	if (val == playerValue) {
		game0remains--;
	}
	score(val == playerValue);
	notify('hideTile(' + line + ',' + col + ')');

	if (game0remains <= 0) {
		finish('You have finished your part!', true);
	}
}
