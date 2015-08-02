"use strict";

var game1remains = 0;

function game1(baseColor, firstColor, secondColor, game1data) {
	if (playerIndex == 1) {
		playerValue = firstColor;
	} else {
		playerValue = secondColor;
	}

	var lines = game1data;
	var list = $('<ul>').addClass('game game1');

	$.each(lines, function(i) {
	    var li = $('<li/>')
	    	// .attr('style', 'height: 35px')
	        .addClass('no-list')
	        .appendTo(list);
	    var subList = $('<ul>').appendTo(li);

	    var colors = lines[i]
	    $.each(colors, function(j) {
	    	var col = colors[j];
	    	if (col == playerValue) {
	    		game1remains++;
	    	}
	    	$('<li onclick="game1Clicked(' + i + ',' + j + ')"/>')	    		
    			.attr('style', 'cursor: pointer; display: inline; color: ' + col)
    			.html("&#x2B24;")
	    		.addClass('no-list circle')
	    		.attr('game-value', col)
	    		.appendTo(subList);
	    });
	});

	return list[0].outerHTML;	
}

function game1Clicked (line, col) {
	if (game1remains <= 0) {
		return;
	}
	hideTile(line, col)
	var $elem = $('ul.game1 li:nth-child(' + (line+1) + ') ul li:nth-child(' + (col+1) + ')');

	var val = $elem.attr('game-value');
	if (val == playerValue) {
		game1remains--;
	}
	score(val == playerValue);
	notify('hideTile(' + line + ',' + col + ')');

	if (game1remains <= 0) {
		finish('You have finished your part!', true);
	}
}