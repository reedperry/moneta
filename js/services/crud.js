services.factory('DB', ['$http', '$log', 'Type', function($http, $log, Type) {

    var CRUD_URL = 'data';

    function storeEvent(event, kind) {
        if (!event) {
            $log.log('Must provide an event object to store.');
            return;
        } else if (!event.amount) {
            $log.log('Event requires an "amount"!')
            return;
        }

        if (kind === Type.EXPENSE || kind === Type.INCOME) {
          event.kind = kind;
        } else {
            $log.log('Invalid kind %s', kind);
            return;
        }

        return $http.post(CRUD_URL, event);
    }

    function getEvents(kind, user) {
      var url = CRUD_URL + '?kind=' + kind;
      if (user) {
        url += '&user=' + user;
      }
      return $http.get(url)
        .success(function(data, status) {
          console.log('Got events: %O', data);
        })
        .error(function(data, status) {
          console.log('HTTP %s: Failed to get events: %O', status, data);
        });
    }

    var db = {

        storeExpense: function(expense) {
          return storeEvent(expense, Type.EXPENSE);
        },

        storeIncome: function(income) {
          return storeEvent(income, Type.INCOME);
        },

        getExpenses: function(user) {
          return getEvents(Type.EXPENSE, user);
        },

        getIncomes: function(user) {
          return getEvents(Type.INCOME, user);
        }
    };

    return db;
}]);
