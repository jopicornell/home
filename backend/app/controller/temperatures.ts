import { Context } from 'aws-lambda';
import { MessageUtil } from '../utils/message';
import { Temperature } from '../model/temperature';
import { injectable } from 'inversify';
import { TemperatureService } from '../service/temperatures';

@injectable()
export class TemperatureController {
  constructor(
    public temperatureService: TemperatureService,
  ) {}

  async create (event: any, context?: Context) {
    console.log('functionName', context.functionName);
    const temperature: Temperature = Temperature.fromJSON(JSON.parse(event.body));

    try {
      const result = await this.temperatureService.create(temperature);

      return MessageUtil.success(result);
    } catch (err) {
      console.error(err);

      return MessageUtil.error();
    }
  }

  async find () {
    try {
      const result = await this.temperatureService.findTemperatures();

      return MessageUtil.success(result);
    } catch (err) {
      console.error(err);

      return MessageUtil.error();
    }
  }
}
