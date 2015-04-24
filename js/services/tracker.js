services.factory('Tracker', function() {

    var expenses = [],
        incomes = [];

    var tracker = {

        /*
         * Add a new expense.
         * Takes a numerical value, or an object containing details of the expense:
         *   amount: number
         *   comment: string
         *   date: date
         */
        addExpense: function(expense) {
            if (!expense) { return; }

            if (typeof expense === 'object') {
                expense.amount = parseFloat(expense.amount);
                if (!expense.date) {
                  expense.date = new Date();
                }
                expenses.push(expense);
            } else if (typeof expense === 'number') {
                var expenseObj = { amount: parseFloat(expense) };
                expenseObj.date = new Date();
                expenses.push(expenseObj);
            }
        },

        /*
         * Add a new income.
         * Takes a numerical value, or an object containing details of the income:
         *   amount: number
         *   comment: string
         *   date: date
         */
        addIncome: function(income) {
            if (!income) { return; }

            if (typeof income === 'object') {
                income.amount = parseFloat(income.amount);
                if (!income.date) {
                  income.date = new Date();
                }
                incomes.push(income);
            } else if (typeof income === 'number') {
                var incomeObj = { amount: parseFloat(income) };
                incomeObj.date = new Date();
                incomes.push(incomeObj);
            }
        },

        getExpenses: function() {
            return expenses;
        },

        getIncomes: function() {
            return incomes;
        },

        findExpense: function(expense) {
            match = _.filter(expenses, _.matches(expense));
            if (_.isEmpty(match)) {
              return null;
            } else {
              return match[0];
            }
        },

        findIncome: function(income) {
            match = _.filter(incomes, _.matches(income));
            if (_.isEmpty(match)) {
              return null;
            } else {
              return match[0];
            }
        },

        totalExpenses: function() {
            return expenses.reduce(function(prev, current) {
                return prev + current.amount;
            }, 0);
        },

        totalIncomes: function() {
            return incomes.reduce(function(prev, current) {
                return prev + current.amount;
            }, 0);
        },

        saveError: function(event) {
            if (!event) {
                return;
            }
            event.saveError = true;
        },

        saveSuccess: function(event) {
            if (!event) {
                return;
            }
            event.saveError = false;
            event.saveSuccess = true;
        }
    };

    return tracker;

});
