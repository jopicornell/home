import { StatusCodes } from 'http-status-codes';

class Result {
  private statusCode: number;
  private data?: any;

  constructor(statusCode: number, data?: any) {
    this.statusCode = statusCode;
    this.data = data;
  }

  /**
   * Serverless: According to the API Gateway specs, the body content must be stringified
   */
  bodyToString () {
    return {
      statusCode: this.statusCode,
      body: JSON.stringify(this.data),
    };
  }
}

export class MessageUtil {
  static success(data: object) {
    const result = new Result(StatusCodes.OK, data);

    return result.bodyToString();
  }

  static error() {
    const result = new Result(StatusCodes.INTERNAL_SERVER_ERROR, null);

    console.log(result.bodyToString());
    return result.bodyToString();
  }
}
