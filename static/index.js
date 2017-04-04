var name, value;

// Set default ajax parameters
$.ajaxSetup({
	contentType: 'application/json',
	dataType: 'json',
	headers: {'X-CSRFToken': Cookies.get('_csrf_token')},
	method: 'POST'
});


// Ban dialog button callbacks
$('.rcon-button-ban-msg').on('click', function(e) {
	// Keep only one message selected
	$('.rcon-button-ban-msg').removeClass('active');
	$(this).toggleClass('active');
	// Enable kick/ban buttons
	$('.kick.button').removeClass('disabled');
	var message_selected = $('.rcon-button-ban-period.active').length;
	if (message_selected > 0) {
		$('.ban.button').removeClass('disabled');
	}
});
$('.rcon-button-ban-period').on('click', function(e) {
	// Keep only one period selected
	$('.rcon-button-ban-period').removeClass('active');
	$(this).toggleClass('active');
	// Enable ban button
	var message_selected = $('.rcon-button-ban-msg.active').length;
	if (message_selected > 0) {
		$('.ban.button').removeClass('disabled');
	}
});

// Open user dialog on user row click
$('.user_action').click(function() {
	name = $(this).data('name');
	value = $(this).data('value');
	$('.value-user').text(name);
	$('.user_dialog').modal('show');
});

// Open map dialog on map row click
$('.map_action').click(function() {
	name = $(this).data('name');
	value = $(this).data('value');
	$('.value-map').text(name);
	$('.map_dialog').modal('show');
});

// Ajax function to request server action
function ajax(url, data) {
	$.ajax({
		async: false,
		data: JSON.stringify(data),
		url: window.location.href + url
	})
	.done(function(data, status, jqxhr) {
		if (data.ok) {
			location.reload(true);
		}
	});
}

// Configure user dialog callbacks
$('.user_dialog').modal({
	onApprove: function(e) {
		if (e.hasClass('kick')) {
			var data = {
				player:  value,
				message: $('.rcon-button-ban-msg.active').data('value')
			};
			ajax('kick', data);
		} else
		if (e.hasClass('ban')) {
			var data = {
				player:  value,
				message: $('.rcon-button-ban-msg.active').data('value'),
				period:  $('.rcon-button-ban-period.active').data('value')
			};
			ajax('ban', data);
		}
		// Reset button states
		$('.rcon-button-ban-msg').removeClass('active');
		$('.rcon-button-ban-period').removeClass('active');
		$('.kick.button').addClass('disabled');
		$('.ban.button').addClass('disabled');
	},
	onDeny: function(e) {
		// Reset button states
		$('.rcon-button-ban-msg').removeClass('active');
		$('.rcon-button-ban-period').removeClass('active');
		$('.kick.button').addClass('disabled');
		$('.ban.button').addClass('disabled');
	}
});

// Configure map dialog callbacks
$('.map_dialog').modal({
	onApprove: function(e) {
		location.reload(true);
		var data = {
			map: value
		};
		ajax('map', data);
	}
});
