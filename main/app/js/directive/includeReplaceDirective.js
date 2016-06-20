angular
  .module('weatherApp')
  .directive('includeReplace', function () {
    return {
      require: 'ngInclude',
      restrict: 'A',
      link: function (scope, el, attrs) {
        el.replaceWith(el.children());
      }
    };
});
