import { Model, PartitionKey } from '@shiftcoders/dynamo-easy';
import { v4 as uuidv4 } from 'uuid';

interface TemperatureRequestJSON {
  id: string;
  temperature: number;
  type: string;
  date: string;
}
console.log(process.env.TEMPERATURE_TABLE);
@Model({
  tableName: process.env.TEMPERATURE_TABLE,
})
class Temperature {
  @PartitionKey()
  public id: string;
  constructor(
    id: string,
    public temperature: number,
    public type: string,
    public date: string,
  ) {
    this.id = id;
  }

  static fromJSON(requestJSON: TemperatureRequestJSON): Temperature {
    let uuid = requestJSON.id;
    if (!uuid) {
      uuid = uuidv4();
    }
    let date = requestJSON.date;
    if (!date) {
      date = new Date().toISOString();
    }
    return new Temperature(uuid, requestJSON.temperature, requestJSON.type, date);
  }
}

export { Temperature, TemperatureRequestJSON };
