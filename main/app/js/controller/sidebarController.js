angular
  .module('weatherApp')
  .controller('SidebarController',
    ['$location', '$scope',
    function($location, $scope) {

      var overlay = $('.overlay');

      $scope.isClosed = true;
      $scope.open = function(path) {
        $scope.toggleNavbar();
        setTimeout(function($location, path) {
          $location.path(path);
          $scope.$apply();
        }, 300, $location, path);
      }
      $scope.toggleNavbar = function toggleNavbar() {
        if ($scope.isClosed) {
          overlay.show();
          $scope.isClosed = false;
        } else {
          overlay.hide();
          $scope.isClosed = true;
        }
      }

}]);
