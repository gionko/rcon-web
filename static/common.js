jQuery.fn.center = function(parent) {
	if (parent) {
		parent = this.parent();
	} else {
		parent = window;
	}
	this.css({
		"position": "absolute",
		"top": ((($(parent).height() - this.outerHeight()) / 2) + $(parent).scrollTop() + "px"),
		"left": ((($(parent).width() - this.outerWidth()) / 2) + $(parent).scrollLeft() + "px")
	});
	return this;
};

jQuery.fn.visible = function() {
	return this.css('visibility', 'visible');
};

jQuery.fn.invisible = function() {
	return this.css('visibility', 'hidden');
};

jQuery.fn.visibility_toggle = function() {
	return this.css('visibility', function(i, visibility) {
		return (visibility == 'visible') ? 'hidden' : 'visible';
	});
};
