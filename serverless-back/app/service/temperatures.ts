import { DynamoStore } from '@shiftcoders/dynamo-easy';
import { injectable } from 'inversify';
import { isAfter, parseISO, sub } from 'date-fns';
import { Temperature } from '../model/temperature';

@injectable()
export class TemperatureService {
    private temperatureStore: DynamoStore<Temperature>;

    constructor() {
      this.temperatureStore = new DynamoStore(Temperature);
    }

    public async create(temperature: Temperature): Promise<Temperature> {
      try {
        await this.temperatureStore.put(temperature).exec();
        return temperature;
      } catch (err) {
        console.error(err);
        throw err;
      }
    }

    public findTemperatures(query: any): Promise<Temperature[]> {
      let temperatureType = query && query.type;
      if (!temperatureType) {
        temperatureType = 'nick-busk';
      }
      return this.temperatureStore.query()
        .index('DateIndex')
        .wherePartitionKey(temperatureType)
        .descending()
        .exec();
    }

    public async shouldNotifyError(): Promise<boolean> {
      const result = await this.temperatureStore.query()
        .index('DateIndex')
        .wherePartitionKey('nick-busk')
        .descending()
        .limit(1)
        .exec();
      if (result.length > 0) {
        console.error('There are no temperatures with nick-busk');
        return;
      }
      const temperature = result[0];
      const date = parseISO(temperature.date);
      const tenMinuteFromNow = sub(new Date(), { minutes: 10 });
      if (!isAfter(date, tenMinuteFromNow)) {
        console.error('NO DATA!');
        return true;
      }
      return false;
    }
}
