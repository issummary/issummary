import {configure} from 'enzyme';
import Adapter from 'enzyme-adapter-react-16';
import moment from 'moment';
import { SpanManager } from '../../src/services/span';
configure({ adapter: new Adapter() });

const dammyMilestone = {
  DueDate: moment([2018, 1, 3]),
  ID: 1,
  IID: 2,
  StartDate: moment([2018, 1, 1]),
  Title: 'dammy milestone',
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

  // describe('#calcEstimateDate', () => {
  //   const spanManager = new SpanManager([dammyMilestone], 2);
  //   it('should return correct date', () => {
  //     const estDate = spanManager.calcEstimateDate(moment([2017, 12, 31]), 5);
  //     console.log(estDate);
  //   });
  // });
});
