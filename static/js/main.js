var demo;

(function($, d, root) {
  Vue.component('countries-table', {
    template: '#countries-table-template',
    props: {
      data: Array,
      columns: Array,
      filterKey: String,
      sourceBurgerCount: Number,
      sourceDollarCost: Number
    },
    computed: {
      filteredData: function () {
        var self = this;
        var filterKey = this.filterKey && this.filterKey.toLowerCase()
        var data = this.data
        var sourceBurgerCount = this.sourceBurgerCount;
        if (filterKey) {
          data = data.filter(function (row) {
            return Object.keys(row).some(function (key) {
              return String(row[key]).toLowerCase().indexOf(filterKey) > -1
            })
          })
        }

        data.forEach(function (d) {
          d.flagClass = `${d.Iso_a2.toLowerCase()} flag`;

          d.targetBurgerCount = self.sourceDollarCost / d.Dollar_price;
          d.localCost = self.sourceBurgerCount * d.Local_price;
          d.targetCost = self.sourceBurgerCount * d.Dollar_price;
          d.burgerCount = Math.ceil(d.targetBurgerCount);
          d.styleObject = { "--frac": d.burgerCount - d.targetBurgerCount };

          if (isNaN(d.targetCost)) {
            d.targetCost = 0;
            d.burgerCount = 0;
            d.styleObject = {};
          }
        });

        return data
      }
    },
  });

  Vue.component('countries-dropdown', {
    template: '#countries-dropdown-template',
    props: {
      country: String,
    },
    mounted: function() {
			this.$nextTick(function() {
        let self = this;
        $(this.$el).dropdown({
          clearable: true,
          onChange: function(value, text, $selectedItem) {
            self.$emit('country-selected', value);
          }
        });
			})
    }
  });

  Vue.component('burgers', {
    template: '#burgers-template',
    props: {
      sourceBurgerCount: Number
    },
    computed: {
      count: function() {
        return Math.ceil(this.sourceBurgerCount);
      }
    }
  });

  app = new Vue({
    el: '#app',
    template: '#app-template',
    data: {
      priceData: d,
      searchQuery: '',
      countryCode: '',
      inputPrice: 0,
    },
    computed: {
      dataForCountry: function() {
        let self = this;
        return self.priceData.find(e => {
          return self.countryCode === e.Iso_a3;
        });
      },
      dollarCost: function() {
        let value = 0.0;
        let tmp = this.dataForCountry;
        if (tmp) value = tmp.Dollar_price * this.sourceBurgerCount;

        return value;
      },
      countrySelected: function() {
        return this.countryCode !== "";
      },
      currencyCode: function() {
        let value = '?';
        let tmp = this.dataForCountry;
        if (tmp) value = tmp.Currency_code;

        return value;
      },
      countryName: function() {
        let value = '?';

        let tmp = this.dataForCountry;
        if (tmp) value = tmp.Name;

        return value;
      },
      sourceBurgerCount: function() {
        if (isNaN(this.inputPrice)) return 0.0;

        let value = 0;
        let tmp = this.dataForCountry;
        if (tmp) value = this.inputPrice / tmp.Local_price;

        return value;
      }
    },
    watch: {
      sourceBurgerCount: function(n) {
        document.querySelector('.burgers').style.
          setProperty('--frac', Math.ceil(n) - n);
      }
    },
    methods: {
      onCountrySelected: function(iso_a3) {
        this.countryCode = iso_a3;
      },
      onInputKeyDown: function(e) {
        let $amount = $(this.$el).find('#price-input-form input[name=amount]');

        switch (e.keyCode) {
          case 13: // Enter
            $amount.blur();
          break;
          case 38: // Up arrow
            this.inputPrice++;
            break;
          case 40: // Down arrow
            this.inputPrice--;
            break;
          default:
            //console.log(e.keyCode)
          break;
        }
      }
    },
    mounted: function() {
      this.$nextTick(function() {
        let $form = $(this.$el.querySelector('#price-input-form'));
        let $amount = $($form.find('input[name=amount]'));

        $form.submit(e => {
          e.preventDefault();
        }).form({
          on: 'change',
          inline: true,
          fields: {
            number: {
              identifier: 'amount',
              rules: [
                {
                  type: 'number',
                  prompt: 'Please enter a valid number'
                }
              ]
            }
          },
          onSuccess: (event, fields) => {
          // noop 
          }
        });
      });
    }
  });
})($, data, document.documentElement);
