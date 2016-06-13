angular.
  module('weatherApp').
  config(['$locationProvider' ,'$routeProvider',
    function config($locationProvider, $routeProvider) {
      $locationProvider.hashPrefix('!');

      $routeProvider
        .when('/overview', {
          templateUrl : 'templates/overview.html',
          controller: 'OverviewController'
        })
        .when('/realtime', {
          templateUrl : 'templates/realtime.html',
          controller: 'RealtimeController'
        })
        .when('/all', {
          templateUrl : 'templates/all.html',
          controller: 'AllController'
        })
        .otherwise('/overview');
    }
  ]);
