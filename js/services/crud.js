services.factory('DB', ['$http', '$log', function($http, $log) {

    var CRUD_URL = 'data';
    var EXPENSE_KIND = 'expense';
    var INCOME_KIND = 'income';

    function storeEvent(event, kind) {
        if (!event) {
            $log.log('Must provide an event object to store.');
            return;
        } else if (!event.amount) {
            $log.log('Event requires an "amount"!')
            return;
        }

        if (kind === EXPENSE_KIND || kind === INCOME_KIND) {
          event.kind = kind;
        } else {
            $log.log('Invalid kind %s', kind);
            return;
        }

        return $http.post(CRUD_URL, event);
    }

    function getEvents(user, kind) {
      if (!user) {
        $log.log('No user to get events.');
        return;
      }
      var url = CRUD_URL + '?user=' + user;
      if (kind) {
        url += '&kind=' + kind;
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
          return storeEvent(expense, EXPENSE_KIND);
        },

        storeIncome: function(income) {
          return storeEvent(income, INCOME_KIND);
        },

        getExpenses: function(user) {
          return getEvents(user, EXPENSE_KIND);
        },

        getIncomes: function(user) {
          return getEvents(user, INCOME_KIND);
        }
    };

    return db;
}]);
