var name, value;

$(function() {
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

	// Configure user dialog callbacks
	$('.user_dialog').modal({
		onApprove: function(e) {
			if (e.hasClass('kick')) {
				var msg = $('.rcon-button-ban-msg.active').data('value');
				// TODO: Add rcon ajax call:
				//       kickid <value> <msg>
				location.reload(true);
			} else
			if (e.hasClass('ban')) {
				var msg = $('.rcon-button-ban-msg.active').data('value');
				var period = $('.rcon-button-ban-period.active').data('value');
				// TODO: Add rcon ajax call:
				//       banid <period> <value>
				//       kickid <value> <msg>
				location.reload(true);
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
			// TODO: Add rcon ajax call:
			//       changelevel <value>
			location.reload(true);
		}
	});
});
