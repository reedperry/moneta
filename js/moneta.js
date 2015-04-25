var filters = angular.module('moneta.filters', []);
var services = angular.module('moneta.services', ['moneta.filters']);
var directives = angular.module('moneta.directives', ['moneta.services', 'moneta.filters']);
var controllers = angular.module('moneta.controllers', ['moneta.services', 'moneta.directives', 'moneta.filters']);

var moneta = angular.module('moneta', ['ngAnimate', 'ngRoute', 'ngTouch', 'moneta.filters',
        'moneta.services', 'moneta.directives', 'moneta.controllers']);

moneta.config(['$compileProvider', '$controllerProvider','$filterProvider', '$locationProvider',
        '$logProvider', '$provide', '$routeProvider',
        function($compileProvider, $controllerProvider, $filterProvider, $locationProvider,
            $logProvider, $provide, $routeProvider) {

    // Available for scripts loaded after the application has been bootstrapped
    moneta.compileProvider = $compileProvider;
    moneta.controllerProvider = $controllerProvider;
    moneta.filterProvider = $filterProvider;
    moneta.provide = $provide;
    moneta.routeProvider = $routeProvider;

    $logProvider.debugEnabled(false);
    $locationProvider.html5Mode(false);

    $routeProvider.when('/', {
        templateUrl: 'views/home.html',
        controller: 'HomeController as home',
    });

    $routeProvider.when('/track', {
        templateUrl: 'views/track.html',
        controller: 'TrackingController as track',
    });


    $routeProvider.when('/reports', {
        templateUrl: 'views/reports.html',
        controller: 'ReportsController as reports',
    });

    $routeProvider.when('/preferences', {
        templateUrl: 'views/prefs.html',
        controller: 'PrefsController as prefs',
    });

    $routeProvider.otherwise({redirectTo: '/'});

}]);

moneta.constant('Type', {
  EXPENSE: 'expense',
  INCOME: 'income',
  EVENT_KIND: 'event'
});
