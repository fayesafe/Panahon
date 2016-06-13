(function($) {
  $(document).ready(function () {
    var trigger = $('.hamburger'),
        overlay = $('.overlay'),
        wrapper = $('#wrapper'),
        links = $('#sidebar-wrapper a'),
        isClosed = false;

      trigger.click(hamburger_cross);

      links.each(function() {
        $(this).click(function() {
          hamburger_cross();
          wrapper.toggleClass('toggled');
        });
      });

      overlay.click(function() {
        hamburger_cross();
        wrapper.toggleClass('toggled');
      });

      function hamburger_cross() {
        if (isClosed == true) {
          overlay.hide();
          trigger.removeClass('is-open');
          trigger.addClass('is-closed');
          isClosed = false;
        } else {
          overlay.show();
          trigger.removeClass('is-closed');
          trigger.addClass('is-open');
          isClosed = true;
        }
    }

    $('[data-toggle="offcanvas"]').click(function () {
        wrapper.toggleClass('toggled');
    });
  });
})(jQuery);
