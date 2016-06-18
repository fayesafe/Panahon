angular
  .module('weatherApp')
  .constant('EVENTS', {
    'DATA_UPDATED': 'data_updated'
  })
  .constant('CHART_OPTIONS', {
    'TEMPERATURE': {
      data: {
        x: 'time',
        rows: [],
        hide: ['humidity', 'sun', 'rain', 'pressure'],
        names: { temperature: 'Temperatur' },
        classes: { temperature: 'temperature' },
        types: { temperature: 'spline' },
        labels: {
          format: {
            temperature: function (v, id, i, j) { return v + '°C'; }
          }
        },
        colors: { temperature: '#ff851b' },
        empty: {
          label: {
            text: "Loading data..."
          }
        }
      },
      legend: {
        hide: ['humidity', 'sun', 'rain', 'pressure']
      },
      axis: {
        x: {
          type: 'timeseries',
          tick: { format: '%H:%M' },
          labels: true
        }
      },
      grid: { y: { lines: [{ value: 0 }] } }
    },
    'HUMIDITY': {
      data: {
        x: 'time',
        rows: [],
        hide: ['temperature', 'sun', 'rain', 'pressure'],
        names: { humidity: 'Luftfeuchtigkeit' },
        classes: { humidity: 'humidity' },
        types: { humidity: 'line' },
        labels: {
          format: {
            humidity: function (v, id, i, j) { return v + '%'; }
          }
        },
        colors: { humidity: '#28b62c' },
        empty: {
          label: {
            text: "Loading data..."
          }
        }
      },
      legend: {
        hide: ['temperature', 'sun', 'rain', 'pressure']
      },
      axis: {
        x: {
          type: 'timeseries',
          tick: { format: '%H:%M' }
        }
      }
    },
    'PRESSURE': {
      data: {
        x: 'time',
        rows: [],
        hide: ['humidity', 'sun', 'rain', 'temperature'],
        names: { pressure: 'Luftdruck' },
        classes: { pressure: 'pressure' },
        type: 'bar',
        colors: { pressure: '#158cba' },
        empty: {
          label: {
            text: "Loading data..."
          }
        }
      },
      legend: {
        hide: ['humidity', 'sun', 'rain', 'temperature']
      },
      axis: {
        x: {
          type: 'timeseries',
          tick: { format: '%H:%M' },
        },
        y: {
          tick: { format: function (d) {
            return d + 'hPa';
          } }
        }
      }
    },
    'RAIN': {
      data: {
        x: 'time',
        rows: [],
        hide: ['humidity', 'sun', 'temperature', 'pressure'],
        names: { rain: 'Regen' },
        classes: { rain: 'rain' },
        types: { rain: 'area-step' },
        colors: { rain: '#158cba' },
        empty: {
          label: {
            text: "Loading data..."
          }
        }
      },
      legend: {
        hide: ['time', 'humidity', 'sun', 'temperature', 'pressure']
      },
      axis: {
        x: {
          type: 'timeseries',
          tick: { format: '%H:%M' },
        },
        y: {
          tick: { format: function (d) {
            return d + '%';
          } }
        }
      }
    },
    'SUN': {
      data: {
        x: 'time',
        rows: [],
        hide: ['humidity', 'rain', 'temperature', 'pressure'],
        names: { sun: 'Sonne' },
        classes: { sun: 'rain' },
        types: { sun: 'spline' },
        colors: { sun: '#ff4136' },
        empty: {
          label: {
            text: "Loading data..."
          }
        }
      },
      legend: {
        hide: ['time', 'humidity', 'rain', 'temperature', 'pressure']
      },
      axis: {
        x: {
          type: 'timeseries',
          tick: { format: '%H:%M' },
        },
        y: {
          inverted: true,
          label: {
            text: 'Helligkeit',
            position: 'outer-middle'
          },
          tick: { format: function (d) { return ''; } }
        }
      }
    },
    'ALL': {
      bindto: '#chart-all',
      size: {
        height: 500
      },
      data: {
        x: 'time',
        rows: [],
        axes: {
          ts: 'x',
          temperature: 'y',
          pressure: 'y2',
          sun: 'y2'
        },
        names: { temperature: 'Temperatur', humidity: 'Luftfeuchtigkeit', pressure: 'Luftdruck', rain: 'Regen', sun: 'Sonne' },
        classes: {  },
        types: { temperature: 'spline', humidity: 'line', pressure: 'bar', rain: 'area-step', sun: 'spline' },
        labels: {
          format: {
            temperature: function (v, id, i, j) { return v + '°C'; },
            humidity: function (v, id, i, j) { return v + '%'; },
            //pressure: function (v, id, i, j) { return v + 'hPa'; },
          }
        },
        colors: { temperature: '#ff851b', humidity: '#28b62c', pressure: '#3498db' },
        empty: { label: { text: "Loading data..." } }
      },
      axis: {
        x: {
          type: 'category',
          labels: true
        },
        y2: {
          show: true,
          tick: { format: function (d) { return d + ' hPa'; } }
        }
      },
      grid: { y: { lines: [{ value: 0 }] } },
      zoom: {
          enabled: true
      },
      subchart: {
          show: true
      }
    }
  })
  .run(function ($rootScope, CHART_OPTIONS) {
    $rootScope.CHART_OPTIONS = CHART_OPTIONS;
  });
