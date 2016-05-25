angular
  .module('weatherApp')
  .constant('EVENTS', {
    'DATA_UPDATED': 'data_updated'
  })
  .constant('CHARTS', {
    'TEMPERATURE': {
      data: [],
      dimensions: {
        time: {
          axis: 'x'
        },
        temp: {
          axis: 'y',
          type: 'spline',
          color: 'orange',
          postfix: '°C',
          name: 'Temperatur'
        }
      },
      chart: {
        axis: {
          x: {
            type: 'timeseries',
            tick: {
              format: '%H:%M:%S'
            }
          }
        },
        grid: {
          x: {
            show: true
          },
          y: {
            show: true,
            lines: [
              { value: 0 }
            ]
          },
        },
        tooltip: {
          format: {
            title: function (d) { return 'Datensatz ' + d; }
          }
        }
      }
    }
  })
