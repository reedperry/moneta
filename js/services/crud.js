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

        return $http.post(CRUD_URL, event).
            success(function(data, status) {
              console.log('Stored %s successully!', kind);
              console.log(data);
            }).
            error(function(data, status) {
              console.log('Failed to store %s - HTTP %s', kind, status);
              console.log(data);
            });
    }

    var db = {

        storeExpense: function(expense) {
          return storeEvent(expense, EXPENSE_KIND);
        },

        storeIncome: function(income) {
          return storeEvent(income, INCOME_KIND);
        }
    };

    return db;
}]);
