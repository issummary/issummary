import Moment from 'moment';
import { Moment as MomentC } from 'moment';
import { DateRange, extendMoment } from 'moment-range';
import { IMilestone } from '../models/milestone';
// import * as koyomi from "koyomi";
const koyomi = require('koyomi');// tslint:disable-line
const moment = extendMoment(Moment);

export class SpanError implements Error {
  public name = 'SpanError';

  constructor(public message: string) {}

  public toString() {
    return this.name + ': ' + this.message;
  }
}

export class Span {// tslint:disable-line
  public range: DateRange;

  constructor(public milestone: IMilestone) {
    if (milestone.StartDate === undefined) {
      throw new SpanError(
        'Invalid Start Date or Due Date of IMilestone: ' + milestone
      );
    }

    if (milestone.DueDate === undefined) {
      throw new SpanError(
        'Invalid Start Date or Due Date of IMilestone: ' + milestone
      );
    }

    this.range = moment.range(milestone.StartDate, milestone.DueDate);
  }
}

export class SpanManager {// tslint:disable-line
  private spans: Span[];
  constructor(milestones: IMilestone[], private defaultParallel: number) {
    this.spans = milestones.map(m => new Span(m));

    if (this.hasDuplicateSpan()) {
      throw new SpanError('duplicated span found');
    }

    this.spans = this.spans.sort((a, b) => {
      if (a.range.start.isBefore(b.range.start)) {
        return 1;
      }
      if (a.range.start.isAfter(b.range.start)) {
        return -1;
      }
      return 0;
    });
  }

  public getContainSpan(date: MomentC): Span | undefined {
    // 重複するSpanがない前提
    return this.spans.find(s => s.range.contains(date));
  }

  public calcEstimateDate(startDate: MomentC, manDays: number): MomentC {
    let remainManDays = manDays;
    const currentDate = startDate;
    while (remainManDays > 0) {
      if (koyomi.isOpen(currentDate.toDate())) {
        const span = this.getContainSpan(currentDate);
        remainManDays -= span ? span.milestone.parallel : this.defaultParallel;
      }
      currentDate.add(1, 'days');
    }
    return currentDate;
  }

  public hasDuplicateSpan(): boolean {
    for (const span of this.spans) {
      const duplicatedSpans = this.spans.filter(s =>
        s.range.contains(span.range)
      );
      if (duplicatedSpans.length > 1) {
        return true;
      }
    }
    return false;
  }
}
