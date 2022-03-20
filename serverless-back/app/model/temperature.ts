import {
  GSIPartitionKey, GSISortKey, Model, PartitionKey,
} from '@shiftcoders/dynamo-easy';
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

    @GSIPartitionKey('DateIndex')
    @GSISortKey('TypeIndex')
    public type: string;

    @GSIPartitionKey('TypeIndex')
    @GSISortKey('DateIndex')
    public date: string;

    constructor(
      id: string,
        public temperature: number,
        type: string,
        date: string,
    ) {
      this.id = id;
      this.type = type;
      this.date = date;
    }

    static fromJSON(requestJSON: TemperatureRequestJSON): Temperature {
      let uuid = requestJSON.id;
      if (!uuid) {
        uuid = uuidv4();
      }
      let { date } = requestJSON;
      if (!date) {
        date = new Date().toISOString();
      }
      return new Temperature(uuid, requestJSON.temperature, requestJSON.type, date);
    }
}

export { Temperature, TemperatureRequestJSON };
