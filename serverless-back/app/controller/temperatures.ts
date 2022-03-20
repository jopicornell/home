import { Context } from 'aws-lambda';
import { injectable } from 'inversify';
import { MessageUtil } from '../utils/message';
import { Temperature } from '../model/temperature';
import { TemperatureService } from '../service/temperatures';

@injectable()
export class TemperatureController {
  constructor(
        public temperatureService: TemperatureService,
  ) {}

  async create(event: any, context?: Context) {
    console.log('functionName', context.functionName);
    try {
      let result : Temperature | Temperature[];
      const jsonRequest = JSON.parse(event.body);
      if (Array.isArray(jsonRequest)) {
        const temperatures: Temperature[] = jsonRequest.map((t) => Temperature.fromJSON(t));
        const resultPromises = temperatures.map((temperature) => this.temperatureService.create(temperature));
        result = await Promise.all(resultPromises);
      } else {
        const temperature: Temperature = Temperature.fromJSON(jsonRequest);
        result = await this.temperatureService.create(temperature);
      }
      console.log('temperature created');
      return MessageUtil.success(result);
    } catch (err) {
      console.error(err);
      return MessageUtil.error();
    }
  }

  async find(event: any, context?: Context) {
    console.log('logging ', event);
    try {
      const result = await this.temperatureService.findTemperatures(event.queryStringParameters);

      return MessageUtil.success(result);
    } catch (err) {
      console.error(err);

      return MessageUtil.error();
    }
  }

  async checkHealth(event: any, context?: Context) {
    console.log('logging ', event);
    try {
      if (await this.temperatureService.shouldNotifyError()) {
        console.error('NOTIFYING!');
      }

      return MessageUtil.success(null);
    } catch (err) {
      console.error(err);
      return MessageUtil.error();
    }
  }
}
