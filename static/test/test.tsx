import { configure, shallow } from 'enzyme';
import * as React from 'react';
import { About } from '../src/components/About';
import Adapter from 'enzyme-adapter-react-16';
configure({ adapter: new Adapter() });

describe('<About />', () => {
  it('have h2', () => {
    const wrapper = shallow(<About />);
    expect(wrapper.contains(<h2>About</h2>)).toBeTruthy();
  });
});
