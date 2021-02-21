import { Temperature } from '../model/temperature';
import { DynamoStore } from '@shiftcoders/dynamo-easy';
import { injectable } from 'inversify';

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

  public findTemperatures(): Promise<Temperature[]> {
    return this.temperatureStore.scan().exec();
  }
}
