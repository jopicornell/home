import { Handler, Context } from 'aws-lambda';
import dotenv from 'dotenv';
import path from 'path';
import 'reflect-metadata';
import { Container } from 'inversify';
import { TemperatureService } from './service/temperatures';
import { TemperatureController } from './controller/temperatures';

const dotenvPath = path.join(__dirname, '../', `config/.env.${process.env.NODE_ENV}`);
dotenv.config({
  path: dotenvPath,
});

if (process.env.NODE_ENV === 'dev') {
  const AWS = require('aws-sdk');
  AWS.config.region = 'localhost';
  AWS.config.dynamodb = {
    endpoint: 'http://localhost:8000',
    accessKeyId: 'DEFAULT_ACCESS_KEY',
    secretAccessKey: 'DEFAULT_SECRET',
  };
}

const container = new Container();

container.bind<TemperatureController>(TemperatureController).to(TemperatureController);
container.bind<TemperatureService>(TemperatureService).to(TemperatureService);

const temperatureController = container.resolve<TemperatureController>(TemperatureController);

export const create: Handler = (event: any, context: Context) => temperatureController.create(event, context);

export const find: Handler = (event: any, context: Context) => temperatureController.find(event, context);
