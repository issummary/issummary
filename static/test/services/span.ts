import { SpanManager } from '../../src/services/span';
import moment from 'moment';

const dammyMilestone = {
  ID: 1,
  IID: 2,
  Title: 'dammy milestone',
  StartDate: moment([2018, 1, 1]),
  DueDate: moment([2018, 1, 3])
};

describe('SpanManager', () => {
  describe('#getContainSpan', () => {
    const spanManager = new SpanManager([dammyMilestone], 2);

    it('should return true if spans contain specified date', () => {
      expect(spanManager.getContainSpan(moment([2018, 1, 1]))).toBeTruthy();
      expect(spanManager.getContainSpan(moment([2018, 1, 2]))).toBeTruthy();
      expect(spanManager.getContainSpan(moment([2018, 1, 3]))).toBeTruthy();
    });

    it('should return false if spans does not contain specified date', () => {
      expect(spanManager.getContainSpan(moment([2017, 12, 31]))).toBeFalsy();
      expect(spanManager.getContainSpan(moment([2018, 1, 4]))).toBeFalsy();
    });
  });

  describe('#calcEstimateDate', () => {
    const spanManager = new SpanManager([dammyMilestone], 2);
    it('should return correct date', () => {
      const estDate = spanManager.calcEstimateDate(moment([2017, 12, 31]), 5);
      console.log(estDate);
    });
  });
});
